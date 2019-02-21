package main

import (
	"bytes"
	"encoding/gob"
	"godistributed-rabbitmq/common"
	"godistributed-rabbitmq/common/dto"
	"godistributed-rabbitmq/storage"
	_ "godistributed-rabbitmq/storage"
	"log"
)

func main() {
	connection, channel := common.GetChannel(common.URL_GUEST)
	defer connection.Close()
	defer channel.Close()

	messages, error := channel.Consume(
		common.PERSISTENCE_QUEUE, "", false, true, false, false, nil)

	if error != nil {
		log.Fatalln("Failed to get access to messages")
	}

	for message := range messages {
		buffer := bytes.NewReader(message.Body)
		decoder := gob.NewDecoder(buffer)
		readout := &dto.Readout{}
		decoder.Decode(readout)

		error := storage.SaveReadout(readout)

		if error != nil {
			log.Printf("Failed to save readings from: %v. Error: %v.", readout.Name, error)
		} else {
			message.Ack(false)
			log.Println("Saved readout from:", readout.Name)
		}
	}
}
