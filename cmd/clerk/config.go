package main

import (
	"flag"
	"fmt"
	"os"

	"protocall/cmd/clerk/application"
	"protocall/cmd/clerk/interfaces"
	"protocall/pkg/logger"
	"protocall/pkg/recognizer"
	"protocall/pkg/s3"

	"gopkg.in/yaml.v2"
)

type config struct {
	Server     interfaces.ServerConfig      `yaml:"server"`
	Logger     logger.LoggerConfig          `yaml:"logger"`
	Recognizer *recognizer.RecognizerConfig `yaml:"recognizer"`
	Storage    *s3.StorageConfig            `yaml:"storage"`
	Sender     *application.SenderConfig    `yaml:"sender"`
}

var (
	configPath = flag.String("f", "", "path to configuration file")
)

func parseConfig() *config {
	flag.Parse()
	if *configPath == "" {
		fmt.Println("need to specify path to config")
		flag.Usage()
		os.Exit(1)
	}
	data, err := os.ReadFile(*configPath)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	config := &config{}
	if err = yaml.Unmarshal(data, config); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	return config
}
