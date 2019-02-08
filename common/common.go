// -----------------------------------------------------------------------------
// Common package used for containing shared code.
// -----------------------------------------------------------------------------
package common

import (
	"fmt"
	"log"

	"github.com/streadway/amqp"
)

// Constants related to the RabbitMQ URLs.
const (
	//URL of the guest user.
	URL_GUEST = "amqp://guest@localhost:5672"
)

// Constants related to the RabbitMQ exchanges names.
const (
	DEFAULT_EXCHANGE   = ""                 // Default direct exchange
	FANOUT_EXCHANGE    = "amq.fanout"       // Default fanout exchange
	DISCOVERY_EXCHANGE = "sensor.discovery" // Custom exchange used for sending discovery requests.
	// TODO: Similar can be done to support sensor unsubscribe action.
)

// Const related to the exchange types.
const (
	DIRECT = "direct"
	FANOUT = "fanout"
	TOPIC  = "topic"
	HEADER = "header"
)

// Const related to the message queues.
const (
	DISCOVERY_QUEUE = "discovery.queue"
)

// -----------------------------------------------------------------------------
// GetChannel - Helper function that returns amqp channele for the given URL.
//
// amqp.Connection - network connection between the application and RabbitMQ.
// amqp.Channel - provides a path fo communication over connection.
// -----------------------------------------------------------------------------
func GetChannel(url string) (*amqp.Connection, *amqp.Channel) {
	connection, error := amqp.Dial(URL_GUEST)
	FailOnError(error, "Failed to connect to RabitMQ")

	channel, error := connection.Channel()
	FailOnError(error, "Failed to optain a channel")

	return connection, channel
}

// -----------------------------------------------------------------------------
// GetQueue - Helper function that returns amqp queue
// for the given queue name and associated channel.
//
// name - name of the requested queue.
// amqp.Channel - provides a path fo communication over connection.
// -----------------------------------------------------------------------------
func GetQueue(name string, channel *amqp.Channel) *amqp.Queue {
	queue, error := channel.QueueDeclare(
		name,  //name
		false, //durable
		false, //autoDelete
		false, //exclusive
		false, //noWait
		nil)   //args
	FailOnError(error, "Failed to declare a queue")
	return &queue
}

// -----------------------------------------------------------------------------
// GetMessageQueue - Helper function that returns message queue whose publishing
// is associated with the default exchange.
//
// amqp.Connection - network connection between the application and RabbitMQ.
// amqp.Channel - provides a path fo communication over connection.
// amqp.Queue - message queue.
// -----------------------------------------------------------------------------
func GetDirectQueue(name string) (*amqp.Connection, *amqp.Channel, *amqp.Queue) {
	connection, channel := GetChannel(URL_GUEST)
	queue := GetQueue(name, channel)
	return connection, channel, queue
}

// -----------------------------------------------------------------------------
// Advertise - Helper function used for publisheshing given name
// to the rest of the system using given advertisement queue.
//
// name - that is advertised to the system.
// amqp.Channel - provides a path fo communication over connection.
// -----------------------------------------------------------------------------
func Advertise(name string, channel *amqp.Channel) {
	message := amqp.Publishing{Body: []byte(name)}
	// Fanout exchange doesn't need queue name to determin where the message goes.
	// It will send the message to every copy of the queue bound to exchange.
	// It's up to the consumer to create message queue.
	channel.Publish(FANOUT_EXCHANGE, "", false, false, message)
}

// -----------------------------------------------------------------------------
// Send - Helper function used for sending slice of data.
//
// data - that is to be sent.
// amqp.Queue - message queue used for sending data.
// amqp.Channel - provides a path fo communication over connection.
// -----------------------------------------------------------------------------
func Send(data []byte, queue *amqp.Queue, channel *amqp.Channel) {
	message := amqp.Publishing{Body: data}
	channel.Publish(DEFAULT_EXCHANGE, queue.Name, false, false, message)
}

// -----------------------------------------------------------------------------
// FailOnError - Checks if the error occured and logs while panicking.
// -----------------------------------------------------------------------------
func FailOnError(error error, message string) {
	if error != nil {
		log.Fatalf("%s: %s", message, error)
		panic(fmt.Sprintf("%s: %s", message, error))
	}
}
