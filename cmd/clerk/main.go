package main

import (
	"context"
	"net/http"

	"protocall/internal/stapler"
	"protocall/internal/stapler/notifier"
	"protocall/internal/translator"
	"protocall/pkg/logger"
	"protocall/pkg/mailer"
	"protocall/pkg/recognizer"
	"protocall/pkg/s3"
	"protocall/pkg/web"
)

func main() {
	cfg := parseConfig()
	logger.NewLogger(&cfg.Logger)

	storage, err := s3.NewStorage(&cfg.Storage)
	if err != nil {
		logger.L.Fatalf("did not connect to s3: %v", err)
	}

	ctx := context.Background()

	rec, err := recognizer.NewRecognizer(ctx, &cfg.Recognizer)
	if err != nil {
		logger.L.Fatalf("did not connect to recognizer: %v", err)
	}

	mux := &http.ServeMux{}
	translator.InitRouter(
		mux,
		&translator.TranslatorHandler{
			App: translator.NewTranslator(rec, storage),
		},
	)
	stapler.InitRouter(
		mux,
		&stapler.StaplerHandler{
			App: stapler.NewStapler(
				storage,
				notifier.NewNotifier(
					mailer.NewMailer(
						&cfg.Mailer,
					),
				),
			),
		},
	)

	logger.L.Infof("Starting server on %v:%v", cfg.Server.Host, cfg.Server.Port)
	if err := web.NewServer(mux, &cfg.Server).Start(); err != nil {
		logger.L.Fatalf("error while running server: %v", err)
	}
}
