package main

import (
	"flag"
	"log"
	"os"

	"protocall/pkg/logger"
	"protocall/pkg/s3"
	"protocall/pkg/webcore"

	"gopkg.in/yaml.v2"
)

type config struct {
	Server  webcore.ServerConfig `yaml:"server"`
	Storage s3.StorageConfig     `yaml:"s3"`
	Logger  logger.LoggerConfig  `yaml:"log"`
	Root    string               `yaml:"root"`
}

var (
	configPath = flag.String("f", "", "path to configuration file")
)

func parseConfig() *config {
	flag.Parse()
	if *configPath == "" {
		flag.Usage()
		log.Fatalf("need to specify path to config")
	}

	data, err := os.ReadFile(*configPath)
	if err != nil {
		log.Fatalf("cannot read configuration: %v", err)
	}

	config := &config{}
	if err = yaml.Unmarshal(data, config); err != nil {
		log.Fatalf("cannot parse configuration: %v", err)
	}

	s3.ApplySecrets(&config.Storage)
	return config
}
