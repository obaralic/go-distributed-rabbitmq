// -----------------------------------------------------------------------------
// Web package used for encapsulating web controllers.
// -----------------------------------------------------------------------------
package controller

import (
	"bytes"
	"encoding/gob"
	"godistributed-rabbitmq/common"
	"godistributed-rabbitmq/common/dto"
	"godistributed-rabbitmq/web/model"
	"log"
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
	"github.com/streadway/amqp"
)

// -----------------------------------------------------------------------------
// message - struct that encapsulates json transport object.
// -----------------------------------------------------------------------------
type message struct {
	Type string      `json:"type"`
	Data interface{} `json:"data"`
}

// -----------------------------------------------------------------------------
// SocketController - struct that encapsulates work with AMQP and Web sockets.
// -----------------------------------------------------------------------------
type SocketController struct {
	connection *amqp.Connection
	channel    *amqp.Channel
	sockets    []*websocket.Conn
	mutex      sync.Mutex
	upgrader   websocket.Upgrader
}

// -----------------------------------------------------------------------------
// NewSocketController - creates new web socket controller.
// -----------------------------------------------------------------------------
func NewSocketController() (controller *SocketController) {
	controller = new(SocketController)
	controller.connection, controller.channel = common.GetChannel(common.URL_GUEST)
	controller.upgrader = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
	}

	go controller.listenForSources()
	go controller.listenForMessages()

	return
}

// -----------------------------------------------------------------------------
// addSocket - adds new web socket to controller.
// -----------------------------------------------------------------------------
func (controller *SocketController) addSocket(socket *websocket.Conn) {
	controller.mutex.Lock()
	log.Print("addSocket")
	controller.sockets = append(controller.sockets, socket)
	controller.mutex.Unlock()
}

// -----------------------------------------------------------------------------
// removeSocket - removes existing web socket from controller.
// -----------------------------------------------------------------------------
func (controller *SocketController) removeSocket(socket *websocket.Conn) {
	controller.mutex.Lock()
	log.Print("removeSocket")
	socket.Close()
	for index := range controller.sockets {
		if controller.sockets[index] == socket {
			controller.sockets = append(controller.sockets[:index], controller.sockets[index+1:]...)

		}
	}

	controller.mutex.Unlock()
}

// -----------------------------------------------------------------------------
// handleMessage -handles incoming http request.
// -----------------------------------------------------------------------------
func (controller *SocketController) handleMessage(writter http.ResponseWriter, request *http.Request) {
	log.Println(request.Header)
	connection, error := controller.upgrader.Upgrade(writter, request, request.Header)

	common.FailOnError(error, "Upgrader error")
	controller.addSocket(connection)

	go controller.listenForDiscoveryRequests(connection)
}

// -----------------------------------------------------------------------------
// listenForSources -
// -----------------------------------------------------------------------------
func (controller *SocketController) listenForSources() {
	queue := common.GetQueue("", controller.channel, true)
	controller.channel.QueueBind(queue.Name, "", common.WEBAPP_SOURCE_EXCHANGE, false, nil)

	messages, _ := controller.channel.Consume(queue.Name, "", true, false, false, false, nil)

	for msg := range messages {
		sensor, error := model.GetSensorByName(string(msg.Body))
		common.FailOnError(error, "Cannot get the sensor by name")
		log.Printf("Soket sensor: %s", sensor.Name)
		controller.sendMessage(message{
			Type: "source",
			Data: sensor,
		})
	}
}

// -----------------------------------------------------------------------------
// listenForMessages -
// -----------------------------------------------------------------------------
func (controller *SocketController) listenForMessages() {
	queue := common.GetQueue("", controller.channel, true)
	controller.channel.QueueBind(
		queue.Name,                      //name string,
		"",                              //key string,
		common.WEBAPP_READINGS_EXCHANGE, //exchange string,
		false,                           //noWait bool,
		nil)                             //args amqp.Table)

	msgs, _ := controller.channel.Consume(
		queue.Name, //queue string,
		"",         //consumer string,
		true,       //autoAck bool,
		false,      //exclusive bool,
		false,      //noLocal bool,
		false,      //noWait bool,
		nil)        //args amqp.Table)

	log.Printf("Web Soket: Waiting for messages")
	for msg := range msgs {
		buffer := bytes.NewBuffer(msg.Body)
		decoder := gob.NewDecoder(buffer)
		readout := dto.Readout{}
		err := decoder.Decode(&readout)

		if err != nil {
			println(err.Error())
		}

		log.Printf("Web Soket: Message received %v", readout)

		controller.sendMessage(message{
			Type: "reading",
			Data: readout,
		})
	}
}

// -----------------------------------------------------------------------------
// listenForDiscoveryRequests -
// -----------------------------------------------------------------------------
func (controller *SocketController) listenForDiscoveryRequests(socket *websocket.Conn) {
	for {
		msg := message{}
		err := socket.ReadJSON(&msg)

		if err != nil {
			log.Printf("Socket Error! %v:", err)
			controller.removeSocket(socket)
			break
		}

		if msg.Type == "discover" {
			controller.channel.Publish(
				"",                               //exchange string,
				common.WEBAPP_DISCOVERY_EXCHANGE, //key string,
				false,                            //mandatory bool,
				false,                            //immediate bool,
				amqp.Publishing{})                //msg amqp.Publishing)
		}
	}
}

// -----------------------------------------------------------------------------
// sendMessage -
// -----------------------------------------------------------------------------
func (controller *SocketController) sendMessage(msg message) {
	socketsToRemove := []*websocket.Conn{}

	log.Printf("Message sources: %v", controller.sockets)
	for _, socket := range controller.sockets {
		err := socket.WriteJSON(msg)

		if err != nil {
			socketsToRemove = append(socketsToRemove, socket)
			common.FailOnError(err, "Cannot Write JSON")
		} else {
			log.Printf("JSON: %v", msg)
		}
	}

	for _, socket := range socketsToRemove {
		controller.removeSocket(socket)
	}
}
