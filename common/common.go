// -----------------------------------------------------------------------------
// Common package used for containing shared code.
// -----------------------------------------------------------------------------
package common

import (
	"fmt"
	"log"

	"github.com/streadway/amqp"
)

const (
	URL_GUEST = "amqp://guest@localhost:5672"
)

// -----------------------------------------------------------------------------
// GetMessageQueue - Helper function that returns message queue whose publishing
// is associated with the default exchange.
//
// amqp.Connection - network connection between the application and RabbitMQ.
// amqp.Channel - provides a path fo communication over connection.
// amqp.Queue - message queue.
// -----------------------------------------------------------------------------
func GetMessageQueue() (*amqp.Connection, *amqp.Channel, *amqp.Queue) {

	connection, error := amqp.Dial(URL_GUEST)
	FailOnError(error, "Failed to connect to RabitMQ")

	channel, error := connection.Channel()
	FailOnError(error, "Failed to optain a channel")

	queue, error := channel.QueueDeclare(
		"Info", // Name of the queue.
		false,  // Will messages be saved to disk when added to the queue.
		false,  // Will messages be deleted when ther is no active consumers.
		false,  // Will queue be accessible only from connection that requests it.
		false,  // Will return exact queue from server or create one if missing.
		nil)    // Will the headers be declared if bound to headers exchange.
	FailOnError(error, "Failed to declare a queue")

	return connection, channel, &queue
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
