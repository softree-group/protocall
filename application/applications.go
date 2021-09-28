package application

import (
	"github.com/CyCoreSystems/ari/v5/client/native"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"protocall/application/app"
	"protocall/application/app/snoopy"
	"protocall/application/applications"
	"protocall/config"
	"protocall/domain/repository"
)

type Applications struct {
	Listener  applications.EventListener
	Connector applications.Connector
	Snoopy    applications.Snoopy
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

	return &Applications{
		Listener:  app.NewListener(ariClient, app.NewHandler(ariClient, reps)),
		Connector: app.NewConnector(ariClient),
		Snoopy:    snoopy.New(),
	}
}
