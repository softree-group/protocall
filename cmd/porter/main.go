package main

import (
	"flag"
	"log"
	"net/http"

	"protocall/pkg/logger"
	"protocall/pkg/s3"
	"protocall/pkg/webcore"
)

var (
	configPath = flag.String("f", "", "path to configuration file")
)

func main() {
	flag.Parse()
	if *configPath == "" {
		flag.Usage()
		log.Fatalf("need to specify path to config")
	}

	cfg, err := config(*configPath)
	if err != nil {
		log.Fatalln(err)
	}

	logger.NewLogger(&cfg.Logger)

	storage, err := s3.NewStorage(&cfg.Storage)
	if err != nil {
		log.Fatalf("cannot connect to s3: %v", err)
	}

	mux := &http.ServeMux{}
	initRouter(mux, &porterHandler{storage, cfg.Root})

	logger.L.Infof("Starting server on %v:%v", cfg.Server.Host, cfg.Server.Port)
	if err = webcore.NewServer(mux, &cfg.Server).Start(); err != nil {
		log.Fatalf("server error: %v", err)
	}
}
