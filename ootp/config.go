package main

import (
	"encoding/json"
	"fmt"
	"os"
	"sort"

	"github.com/google/go-cmp/cmp"
)

// Config is the top-level structure for the configuration JSON.
type Config struct {
	ConfigNumber int      `json:"config_number"`
	Modules      []Module `json:"modules"`
}

// Module represents each item in the modules array.
type Module struct {
	ContainerName  string            `json:"container_name"`
	ContainerImage string            `json:"container_image"`
	RestartPolicy  string            `json:"restart_policy"`
	PullPolicy     string            `json:"pull_policy"`
	StartupOrder   int               `json:"startup_order"`
	Command        string            `json:"command"`
	EnvVariables   map[string]string `json:"env_variables"`
}

func check(e error) {
	if e != nil {
		fmt.Println(e)
	}
}

func sortbyStartupOrder(modules []Module) []Module {
	sort.Slice(modules, func(i, j int) bool {
		return modules[i].StartupOrder < modules[j].StartupOrder
	})
	return modules
}

// define a function to read the config.json and return a Config struct
func ReadConfig(file string) Config {

	jsonStr, err := os.ReadFile("config.json")
	check(err)

	var config Config
	check(json.Unmarshal([]byte(string(jsonStr)), &config))
	return config
}

func CompareConfigs(config1, config2 Config) []Module {
	var changedModules []Module

	for _, module1 := range config1.Modules {
		found := false
		for _, module2 := range config2.Modules {
			if module1.ContainerName == module2.ContainerName {
				found = true
				if !cmp.Equal(module1, module2) {
					changedModules = append(changedModules, module1)
				}
				break
			}
		}
		if !found {
			changedModules = append(changedModules, module1)
		}
	}

	return changedModules
}

func needRestart(firstModuleToStart Module, modules []Module) []Module {
	var modulesToRestart []Module

	for _, module := range modules {
		if module.StartupOrder >= firstModuleToStart.StartupOrder {
			modulesToRestart = append(modulesToRestart, module)
		}
	}

	return modulesToRestart
}
