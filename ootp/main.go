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

	var testModule = config.Modules[0]

	// Example usage
	docker, err := NewDocker()
	if err != nil {
		panic(err)
	}

	err = docker.PullImage(testModule)
	if err != nil {
		fmt.Println("Error pulling image:", err)
		return
	}

	err = docker.RunContainer(testModule)
	if err != nil {
		fmt.Println("Error running container:", err)
		return
	}

	containers, err := docker.ListContainers()
	if err != nil {
		fmt.Println("Error listing containers:", err)
		return
	}

	for _, container := range containers {
		fmt.Println(container.ID)
	}

	err = docker.StopContainer(testModule.ContainerName)
	if err != nil {
		fmt.Println("Error stopping container:", err)
		return
	}

	err = docker.RemoveContainer(testModule.ContainerName)
	if err != nil {
		fmt.Println("Error removing container:", err)
		return
	}

	<-make(chan struct{})
}
