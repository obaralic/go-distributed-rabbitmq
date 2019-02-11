package main

import (
	"fmt"
	"godistributed-rabbitmq/coordinator"
)

var consumer *coordinator.StorageConsumer

func main() {
	fmt.Println("Starting sensor listener...")
	aggregator := coordinator.NewAggregator()
	consumer = coordinator.NewStorageConsumer(aggregator)
	listener := coordinator.NewListener(aggregator)
	defer listener.Stop()

	go listener.Start()

	var input string
	fmt.Scanln(&input)
}
