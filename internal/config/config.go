package config

import (
	"io"
	"log"
	"os"

	"gopkg.in/yaml.v3"
)

type GPIO struct {
	Name     string `yaml:"name"`
	Location string `yaml:"location"`
}

type GPIOConfig struct {
	GPIOs []GPIO `yaml:"gpio"`
}

func GetConfig() GPIOConfig {
	file, err := os.Open("/etc/gpio.yml")
	if err != nil {
		log.Fatalf("error: %v", err)
	}
	defer file.Close()

	data, err := io.ReadAll(file)
	if err != nil {
		log.Fatalf("error: %v", err)
	}

	var config GPIOConfig
	err = yaml.Unmarshal(data, &config)
	if err != nil {
		log.Fatalf("error: %v", err)
	}

	return config
}
