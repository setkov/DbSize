package main

import (
	"os"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Connections []string
	WebUI       struct {
		Port       int
		ChromePath string
	}
}

func NewConfig(fileName string) (*Config, error) {
	buf, err := os.ReadFile(fileName)
	if err != nil {
		return nil, err
	}

	config := &Config{}
	err = yaml.Unmarshal(buf, config)
	if err != nil {
		return nil, err
	}

	return config, nil
}
