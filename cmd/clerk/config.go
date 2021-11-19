package main

import (
	"flag"
	"log"
	"os"

	"protocall/pkg/logger"
	"protocall/pkg/mailer"
	"protocall/pkg/recognizer"
	"protocall/pkg/s3"
	"protocall/pkg/web"

	"gopkg.in/yaml.v2"
)

type config struct {
	Server     web.ServerConfig            `yaml:"server"`
	Logger     logger.LoggerConfig         `yaml:"log"`
	Recognizer recognizer.RecognizerConfig `yaml:"recognizer"`
	Storage    s3.StorageConfig            `yaml:"s3"`
	Mailer     mailer.MailerConfig         `yaml:"smtp"`
}

var (
	configPath = flag.String("f", "", "path to configuration file")
)

func parseConfig() *config {
	flag.Parse()
	if *configPath == "" {
		flag.Usage()
		log.Fatalln("need to specify path to config")
	}
	data, err := os.ReadFile(*configPath)
	if err != nil {
		log.Fatalln(err)
	}

	config := &config{}
	if err = yaml.Unmarshal(data, config); err != nil {
		log.Fatalln(err)
	}

	s3.ApplySecrets(&config.Storage)
	mailer.ApplySecrets(&config.Mailer)
	return config
}
