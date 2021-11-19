package main

import (
	"flag"
	"fmt"
	"os"

	"protocall/pkg/logger"
	"protocall/pkg/mailer"
	"protocall/pkg/recognizer"
	"protocall/pkg/s3"
	"protocall/pkg/web"

	"gopkg.in/yaml.v2"
)

type config struct {
	Server     *web.ServerConfig            `yaml:"server"`
	Logger     *logger.LoggerConfig         `yaml:"log"`
	Recognizer *recognizer.RecognizerConfig `yaml:"recognizer"`
	Storage    *s3.StorageConfig            `yaml:"s3"`
	Mailer     *mailer.MailerConfig         `yaml:"smtp"`
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

	s3.ApplySecrets(config.Storage)
	mailer.ApplySecrets(config.Mailer)
	return config
}
