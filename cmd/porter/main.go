package main

import (
	"log"

	"protocall/internal/porter"
	"protocall/pkg/logger"
	"protocall/pkg/s3"
	"protocall/pkg/web"
)

func main() {
	cfg := parseConfig()
	logger.NewLogger(cfg.Logger)

	storage, err := s3.NewStorage(cfg.Storage)
	if err != nil {
		log.Fatalf("cannot connect to s3: %v", err)
	}

	srv := web.NewServer(cfg.Server, porter.NewRouter(&porter.PorterHandler{
		Storage: storage,
		Root:    cfg.Root,
	}))

	if err = srv.Start(); err != nil {
		log.Fatalf("server error: %v", err)
	}
}
