package main

import (
	"fmt"
	"log"
)

func main() {
	var config = ReadConfig("config.json")

	configUpdates := make(chan string)

	cli := make(chan string)

	docker, err := NewDocker()

	orchestrator := NewOrchestrator(&config, docker)

	commuication, err := NewCommunication("/tmp/ootp.sock", cli)

	listener, err := NewFileListener("config.json", configUpdates)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Starting OOTP")

	orchestrator.Start()

	go listener.Watch()

	go orchestrator.Watch(configUpdates)

	go orchestrator.Listen(cli)

	go commuication.Listen()

	<-make(chan struct{})
}
