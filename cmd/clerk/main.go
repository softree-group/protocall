package main

import (
	"context"
	"protocall/internal/clerk/translator"
	"protocall/pkg/recognizer"
	"protocall/pkg/web"
)

func main() {
	cfg := parseConfig()
	logger.NewLogger(cfg.Logger)

	storage, err := s3.NewStorage(cfg.Storage)
	if err != nil {
		logger.L.Fatalf("did not connect to s3: %v", err)
	}

	ctx := context.Background()

	rec, err := recognizer.NewRecognizer(ctx, cfg.Recognizer)
	if err != nil {
		logger.L.Fatalf("did not connect to recognizer: %v", err)
	}

	if err := server.NewServer(
		cfg.Server,
		clerk.NewRouter(clerk.Handler{
			.NewSender(mailer.NewMailer(cfg.Mailer)),
			application.NewStapler(storage),
			application.NewTranslator(rec, storage)
		})
	)

	if err := web.NewServer(
		cfg.Server,
		clerk.NewRouter(clerk.ClerkHandler{
			Translator: translator.NewTranslator()
		})
	).Start(); err != nil {
		logger.L.Fatalf("error while running server: %v", err)
	}
}
