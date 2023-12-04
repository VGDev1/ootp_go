package main

import (
	"fmt"
)

type Orchestrator struct {
	Config *Config
	Docker *Docker
}

func NewOrchestrator(config *Config, docker *Docker) *Orchestrator {
	return &Orchestrator{Config: config, Docker: docker}
}

func (o *Orchestrator) Start() {
	o.startModules(o.Config.Modules)
}

func (o *Orchestrator) Watch(configUpdates chan string) {
	for {
		select {
		case update := <-configUpdates:
			fmt.Println(update)
			configNew, err := ReadConfig("config.json")
			if err != nil {
				fmt.Println("Error reading config:", err)
				continue
			}
			configOld := *o.Config
			changedModules := CompareConfigs(configOld, *configNew)

			if len(changedModules) == 0 {
				fmt.Println("No modules changed")
				continue
			}
			changedModules = sortbyStartupOrder(changedModules)
			firstModuleToStart := changedModules[0]

			modulesToRestart := needRestart(firstModuleToStart, configNew.Modules)

			o.startModules(modulesToRestart)
			*o.Config = *configNew
		}
	}
}

func (o *Orchestrator) Listen(cli chan string) {
	for {
		select {
		case command := <-cli:
			fmt.Println("Received command:", command)
			switch command {
			case "restart":
				fmt.Println("Restarting")
				o.startModules(o.Config.Modules)
			case "list":
				fmt.Println("Listing modules")
				for _, module := range o.Config.Modules {
					fmt.Println(module.ContainerName)
				}
			}
		}
	}
}

func (o *Orchestrator) startModules(modules []Module) {

	var err error

	modules = sortbyStartupOrder(modules)

	for _, module := range modules {

		fmt.Println("Starting module", module.ContainerName)

		if err := o.Docker.StopContainer(module.ContainerName); err != nil {
			fmt.Println("Error stopping container:", err)
		}

		if err := o.Docker.RemoveContainer(module.ContainerName); err != nil {
			fmt.Println("Error removing container:", err)
		}

		switch module.PullPolicy {
		case "never":
			// Do nothing
		case "ifnotpresent":
			containers, err := o.Docker.ListContainers()
			if err != nil {
				fmt.Println("Error listing containers:", err)
				return
			}
			for _, container := range containers {
				for _, name := range container.Names {
					if name == "/"+module.ContainerName {
						fmt.Println("Container", module.ContainerName, "already running")
						continue
					}
				}
			}
		case "always":
			err = o.Docker.PullImage(module)
			if err != nil {
				fmt.Println("Error pulling image:", err)
				return
			}
		default:
			fmt.Println("Unknown pull policy:", module.PullPolicy)
			return
		}

		err = o.Docker.RunContainer(module)
		if err != nil {
			fmt.Println("Error running container:", err)
			return
		}
	}
}
