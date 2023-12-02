package main

import (
	"fmt"
	"log"
)

func main() {
	var config = ReadConfig("config.json")

	messages := make(chan string)

	listener, err := NewFileListener("config.json", messages)
	if err != nil {
		log.Fatal(err)
	}

	go listener.Watch()

	go func() {
		for {
			select {
			case message := <-messages:
				fmt.Println(message)
				configNew := ReadConfig("config.json")
				changedModules := CompareConfigs(config, configNew)
				fmt.Println(changedModules)
				config = configNew
			}
		}
	}()

	<-make(chan struct{})
}
