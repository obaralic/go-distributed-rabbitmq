// -----------------------------------------------------------------------------
// Main package used for initial tampering with the RabbitMQ server
// -----------------------------------------------------------------------------
package main

import (
	"fmt"
	"godistributed-rabbitmq/common"
	"math/rand"
	"strconv"
	"time"

	"github.com/streadway/amqp"
)

func main() {
	go serverMock()

	var input string
	fmt.Scanln(&input)
}

func serverMock() {
	connection, channel, queue := common.GetMessageQueue()
	defer connection.Close()
	defer channel.Close()

	for {
		index := rand.Intn(1000)

		// Publishing is the struct for encapsulating sending message.
		message := amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte("Health check: "),
		}
		message.Body = append(message.Body, []byte(strconv.Itoa(index))...)

		channel.Publish(
			"",         // Name of the exchange that is used. Empty for the default one.
			queue.Name, // Routing key for resolving which queue should get the message.
			false,      // Make sure that message is delivered.
			false,      // Makse sure when message is delivered.
			message)

		time.Sleep(time.Second)
	}
}
