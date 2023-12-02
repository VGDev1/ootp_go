package main

import (
	"encoding/json"
	"os"

	"github.com/google/go-cmp/cmp"
)

// Config is the top-level structure for the configuration JSON.
type Config struct {
	ConfigNumber int      `json:"config_number"`
	Modules      []Module `json:"modules"`
}

// Module represents each item in the modules array.
type Module struct {
	ModuleName     string            `json:"module_name"`
	ContainerImage string            `json:"container_image"`
	RestartPolicy  string            `json:"restart_policy"`
	PullPolicy     string            `json:"pull_policy"`
	StartupOrder   int               `json:"startup_order"`
	Command        string            `json:"command"`
	EnvVariables   map[string]string `json:"env_variables"`
}

func check(e error) {
	if e != nil {
		panic(e)
	}
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
			if module1.ModuleName == module2.ModuleName {
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
