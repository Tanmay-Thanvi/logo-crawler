package config

import (
	"log"
	"os"

	"gopkg.in/yaml.v3"
)

type Preferences struct {
	Preferred struct {
		MinWidth  int `yaml:"min_width"`
		MinHeight int `yaml:"min_height"`
	} `yaml:"preferred"`
}

func LoadConfig(path string) Preferences {
	var cfg Preferences
	data, err := os.ReadFile(path)
	if err != nil {
		log.Fatalf("Failed to read config file: %v", err)
	}
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		log.Fatalf("Failed to parse config.yaml: %v", err)
	}
	return cfg
}
