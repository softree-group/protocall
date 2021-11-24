package main

import (
	"net/http"

	"protocall/internal/notifier"
	"protocall/internal/stapler"
	"protocall/internal/translator"
	"protocall/pkg/connector"
	"protocall/pkg/logger"
	"protocall/pkg/mailer"
	"protocall/pkg/s3"
	"protocall/pkg/webcore"
	"protocall/pkg/yastt"
)

func main() {
	cfg := parseConfig()
	logger.NewLogger(&cfg.Logger)

	storage, err := s3.NewStorage(&cfg.Storage)
	if err != nil {
		logger.L.Fatalf("did not connect to s3: %v", err)
	}

	mux := &http.ServeMux{}
	translator.InitRouter(
		mux,
		&translator.TranslatorHandler{
			App: translator.NewTranslator(
				yastt.NewYastt(http.DefaultClient, &cfg.Recognizer),
				storage,
				connector.NewConnectorCLient(
					http.DefaultClient,
					&cfg.Connector,
				),
			),
		},
	)
	stapler.InitRouter(
		mux,
		&stapler.StaplerHandler{
			stapler.NewStapler(storage),
			notifier.NewNotifier(mailer.NewMailer(&cfg.Mailer)),
		},
	)

	logger.L.Infof("Starting server on %v:%v", cfg.Server.Host, cfg.Server.Port)
	if err := webcore.NewServer(mux, &cfg.Server).Start(); err != nil {
		logger.L.Fatalf("error while running server: %v", err)
	}
}
