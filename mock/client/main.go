// -----------------------------------------------------------------------------
// Main package used for initial tampering with the RabbitMQ client
// -----------------------------------------------------------------------------
package main

import (
	"fmt"
	"godistributed-rabbitmq/common"
	"log"
)

func main() {
	go clientMock()

	var input string
	fmt.Scanln(&input)
}

func clientMock() {
	connection, channel, queue := common.GetMessageQueue("Info")
	defer connection.Close()
	defer channel.Close()

	messages, error := channel.Consume(
		queue.Name, // Name of the exchange that is used. Empty for the default one.
		"",         // Distinct identifier of client. Empty to assigne one by Rabbit.
		true,       // Automatic acknowledge successfull receiving.
		false,      // Make sure that this is the only one client.
		false,      // Make sure that consumers are on the same connection as sender.
		false,      // makse sure it waits.
		nil)        // N/A
	common.FailOnError(error, "Failed to register a consumer")

	for message := range messages {
		log.Printf("Received message with body: %s", message.Body)
	}
}
