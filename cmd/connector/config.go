package main

import (
	"os"

	"protocall/internal/bridge"
	"protocall/internal/operator"
	"protocall/pkg/centrifugo"
	"protocall/pkg/logger"

	"gopkg.in/yaml.v2"
)

type Config struct {
	Centrifugo centrifugo.Config `yaml:"centrifugo"`
	Logger     logger.Config     `yaml:"log"`
	Upstream   bridge.Config     `yaml:"ari"`
	Operator   operator.Config   `yaml:"ami"`
	Server     struct {
		APIKey         string   `yaml:"key"`
		Host           string   `yaml:"host"`
		Port           string   `yaml:"port"`
		Domain         string   `yaml:"domain"`
		AllowedDomains []string `yaml:"allowedDomains"`
	} `yaml:"server"`
}

func config(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	config := &Config{}
	if err = yaml.Unmarshal(data, config); err != nil {
		return nil, err
	}
	return config, nil
}
