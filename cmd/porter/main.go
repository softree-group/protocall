package main

import (
	"log"
	"net/http"

	"protocall/pkg/logger"
	"protocall/pkg/s3"
	"protocall/pkg/webcore"
)

func main() {
	cfg := parseConfig()
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
