package application

import (
	"protocall/application/app"
	"protocall/application/applications"
	"protocall/domain/repository"
	"protocall/internal/config"

	"github.com/CyCoreSystems/ari/v5/client/native"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

type Applications struct {
	Listener        applications.EventListener
	Snoopy          applications.Snoopy
	User            applications.User
	AsteriskAccount applications.AsteriskAccount
	Conference      applications.Conference
	Connector       applications.Connector
}

func New(reps repository.Repositories) *Applications {
	logrus.Info("start connect")
	ariClient, err := native.Connect(&native.Options{
		Application:  viper.GetString(config.ARIApplication),
		URL:          viper.GetString(config.ARIUrl),
		WebsocketURL: viper.GetString(config.ARIWebsocketURL),
		Username:     viper.GetString(config.ARIUser),
		Password:     viper.GetString(config.ARIPassword),
	})
	logrus.Info("end connect")
	if err != nil {
		logrus.Fatal("cannot connect to asterisk: ", err)
	}

	connector := app.NewConnector(ariClient, reps)

	return &Applications{
		Listener:        app.NewListener(reps, ariClient, app.NewHandler(ariClient, reps, connector)),
		Snoopy:          app.NewSnoopy(),
		Conference:      app.NewConference(reps, ariClient),
		AsteriskAccount: app.NewAsteriskAccount(reps, viper.GetString(config.ARIAccountsFile)),
		User:            app.NewUser(reps),
		Connector:       connector,
	}
}
