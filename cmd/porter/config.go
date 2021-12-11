package main

import (
	"os"

	"protocall/pkg/logger"
	"protocall/pkg/s3"
	"protocall/pkg/webcore"

	"gopkg.in/yaml.v2"
)

type Config struct {
	Server  webcore.ServerConfig `yaml:"server"`
	Storage s3.StorageConfig     `yaml:"s3"`
	Logger  logger.Config        `yaml:"log"`
	Root    string               `yaml:"root"`
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
