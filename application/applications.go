package application

import (
	"protocall/application/app"
	"protocall/application/app/snoopy"
	"protocall/application/applications"
	"protocall/config"
	"protocall/domain/repository"
	"protocall/infrastructure/storage"

	"github.com/CyCoreSystems/ari/v5/client/native"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

type Applications struct {
	Listener applications.EventListener
	Snoopy   applications.Snoopy
	Voice    app.Voice
}

func New(reps *repository.Repositories) *Applications {
	logrus.Info("start connect")
	ariClient, err := native.Connect(&native.Options{
		Application:  viper.GetString(config.ARIApplication),
		URL:          viper.GetString(config.ARIUrl),
		WebsocketURL: viper.GetString(config.ARIWebsocketUrl),
		Username:     viper.GetString(config.ARIUser),
		Password:     viper.GetString(config.ARIPassword),
	})
	logrus.Info("end connect")
	if err != nil {
		logrus.Fatal("cannot connect to asterisk: ", err)
	}

	storage, err := storage.NewStorage(&storage.Config{

	})
	if err != nil {
		logrus.Fatal("cannot connect to s3: ", err)
	}

	r, err := recognizer.New

	return &Applications{
		Listener: app.NewListener(ariClient, app.NewHandler(ariClient, reps, app.NewConnector(ariClient, reps.Bridge))),
		Snoopy:   snoopy.New(),
		Voice:    app.NewVoice(),
	}
}
