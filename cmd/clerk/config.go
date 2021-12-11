package main

import (
	"os"

	"protocall/pkg/connector"
	"protocall/pkg/logger"
	"protocall/pkg/mailer"
	"protocall/pkg/s3"
	"protocall/pkg/webcore"
	"protocall/pkg/yastt"

	"gopkg.in/yaml.v2"
)

type Config struct {
	Server     webcore.ServerConfig            `yaml:"server"`
	Logger     logger.Config                   `yaml:"log"`
	Recognizer yastt.YasttConfig               `yaml:"yastt"`
	Connector  connector.ConnectorClientConfig `yaml:"connector"`
	Storage    s3.StorageConfig                `yaml:"s3"`
	Mailer     mailer.MailerConfig             `yaml:"smtp"`
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
