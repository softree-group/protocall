package main

import (
	"context"
	"os"

	"protocall/cmd/clerk/application"
	"protocall/cmd/clerk/interfaces"
	"protocall/pkg/logger"
	"protocall/pkg/recognizer"
	"protocall/pkg/s3"
)

func main() {
	config := parseConfig()
	logger.NewLogger(&config.Logger)

	storage, err := s3.NewStorage(&s3.StorageConfig{
		Bucket:    config.Storage.Bucket,
		Endpoint:  config.Storage.Endpoint,
		AccessKey: os.Getenv("ACCESS_KEY"),
		SecretKey: os.Getenv("SECRET_KEY"),
	})
	if err != nil {
		logger.L.Fatalf("did not connect to s3: %v", err)
	}

	ctx := context.Background()

	rec, err := recognizer.NewRecognizer(ctx, config.Recognizer)
	if err != nil {
		logger.L.Fatalf("did not connect to recognizer: %v", err)
	}

	if err := interfaces.NewServer(
		&config.Server,
		interfaces.NewRouter(
			interfaces.NewApplication(
				application.NewSender(&application.SenderConfig{
					Host:     config.Sender.Host,
					Port:     config.Sender.Port,
					Username: config.Sender.Username,
					Password: os.Getenv("EMAIL_KEY"),
				}),
				application.NewStapler(storage),
				application.NewTranslator(rec, storage),
			),
		),
	).Start(); err != nil {
		logger.L.Fatalf("error while running server: %v", err)
	}
}
