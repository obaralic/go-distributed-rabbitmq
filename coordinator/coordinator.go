// -----------------------------------------------------------------------------
// Coordinator package used for defining coordinator as a queue listener.
// -----------------------------------------------------------------------------
package coordinator

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"godistributed-rabbitmq/common"
	"godistributed-rabbitmq/common/dto"

	"github.com/streadway/amqp"
)

// -----------------------------------------------------------------------------
// SensorListener - Struct that contains the logic for discovering
// sensor data queues, receiving messages and translating them into events.
// -----------------------------------------------------------------------------
type SensorListener struct {
	// Connection to RabbitMQ.
	connection *amqp.Connection

	// Channel created with the connection.
	channel *amqp.Channel

	// Subscribed sensors with their delivery channels.
	sources map[string]<-chan amqp.Delivery
}

// -----------------------------------------------------------------------------
// NewListener - Creates new sensor listener.
// -----------------------------------------------------------------------------
func NewListener() *SensorListener {
	listener := SensorListener{
		sources: make(map[string]<-chan amqp.Delivery),
	}

	listener.connection, listener.channel = common.GetChannel(common.URL_GUEST)
	return &listener
}

// -----------------------------------------------------------------------------
// Start - Method used for starting sensor observation process.
// SensorListener will receive advertisement messages when the new sensor
// gets pluged into the system.
// -----------------------------------------------------------------------------
func (listener *SensorListener) Start() {
	// Passing "" will result with unique queue name creation by RabbitMQ.
	queue := common.GetQueue("", listener.channel)

	// Rebind queue from default exchange to the fanout.
	listener.channel.QueueBind(queue.Name, "", common.FANOUT_EXCHANGE, false, nil)

	// Receive sensor advertisement message when the new sensor is up and running.
	advertisements, _ := listener.channel.Consume(
		queue.Name, "", true, false, false, false, nil)

	for advertisement := range advertisements {

		sensorName := string(advertisement.Body)

		sensor, _ := listener.channel.Consume(
			sensorName, "", true, false, false, false, nil)

		if listener.sources[sensorName] == nil {
			listener.sources[sensorName] = sensor

			// Launch goroutine for observing incoming sensor readouts.
			go listener.observe(sensor)
		}
	}
}

// -----------------------------------------------------------------------------
// Stop - Method used for stoping sensor observation process.
// SensorListener will receive shutdown messages when the sensor
// gets pluged out of the system.
// -----------------------------------------------------------------------------
func (listener *SensorListener) Stop() {
	defer listener.channel.Close()
	defer listener.connection.Close()
}

// -----------------------------------------------------------------------------
// observe - Method used for observing incoming messages
// received from the subscribed sensor channel.
//
// sensor - incoming channel of amqp.Delivery containing sensor messages.
// -----------------------------------------------------------------------------
func (listener *SensorListener) observe(sensor <-chan amqp.Delivery) {
	for data := range sensor {
		payload := data.Body
		reader := bytes.NewReader(payload)
		decoder := gob.NewDecoder(reader)

		readout := new(dto.Readout)
		decoder.Decode(readout)

		fmt.Printf("Received readout: %v\n", readout)
	}
}
