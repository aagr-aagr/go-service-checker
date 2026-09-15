package main

import (
	"encoding/json"
	"os"
)

type ServiceConfig struct {
	Name string `json:"name"`
	Host string `json:"host"`
	Port int    `json:"port"`
}

func readConfig(path string) ([]ServiceConfig, error) {

	file, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var services []ServiceConfig
	err = json.Unmarshal(file, &services)
	if err != nil {
		return nil, err
	}

	return services, nil

}
