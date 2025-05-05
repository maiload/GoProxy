package config

import (
	"gopkg.in/yaml.v3"
	"os"
)

type Config struct {
	Server struct {
		Port string `yaml:"port"`
		SSL  struct {
			Port string `yaml:"port"`
			Cert string `yaml:"cert"`
			Key  string `yaml:"key"`
		} `yaml:"ssl"`
	} `yaml:"server"`
	Routes []struct {
		Path   string `yaml:"path"`
		Target string `yaml:"target"`
	} `yaml:"routes"`
}

func LoadConfig(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var config Config
	if err := yaml.Unmarshal(data, &config); err != nil {
		return nil, err
	}
	return &config, nil
}
