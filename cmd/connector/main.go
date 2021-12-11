package main

import (
	"flag"
	"log"
	"os"

	"protocall/pkg/logger"
)

var (
	configPath = flag.String("f", "", "path to configuration file")
)

func main() {
	flag.Parse()
	if *configPath == "" {
		flag.Usage()
		os.Exit(1)
	}
	cfg, err := config(*configPath)
	if err != nil {
		log.Fatalln(err)
	}
	logger.NewLogger(&cfg.Logger)
	logger.L.Info(app())
}
