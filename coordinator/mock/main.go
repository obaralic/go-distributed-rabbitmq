package main

import (
	"fmt"
	"godistributed-rabbitmq/coordinator"
)

func main() {
	fmt.Println("Starting sensor listener...")
	listener := coordinator.NewListener()
	defer listener.Stop()

	go listener.Start()

	var input string
	fmt.Scanln(&input)
}
