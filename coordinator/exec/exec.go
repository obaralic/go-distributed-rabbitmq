package main

import (
	"fmt"
	"godistributed-rabbitmq/coordinator"
	_ "godistributed-rabbitmq/storage"
	"log"
)

var dbConsumer *coordinator.StorageConsumer
var webConsumer *coordinator.WebappConsumer

func main() {
	log.Println("Starting sensor listener...")
	aggregator := coordinator.NewAggregator()
	dbConsumer = coordinator.NewStorageConsumer(aggregator)
	webConsumer = coordinator.NewWebappConsumer(aggregator)
	listener := coordinator.NewListener(aggregator)
	defer listener.Stop()

	go listener.Start()

	var input string
	fmt.Scanln(&input)
}
