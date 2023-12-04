package main

import (
	"fmt"
	"log"
)

func main() {
	var config = ReadConfig("config.json")

	docker, err := NewDocker()

	orchestrator := NewOrchestrator(&config, docker)

	configUpdates := make(chan string)

	listener, err := NewFileListener("config.json", configUpdates)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Starting OOTP")

	orchestrator.Start()

	go listener.Watch()

	go orchestrator.Watch(configUpdates)

	<-make(chan struct{})
}
