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
	Listener        applications.EventListener
	Snoopy          applications.Snoopy
	User            applications.User
	AsteriskAccount applications.AsteriskAccount
	Conference      applications.Conference
	Connector       applications.Connector
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

	connector := app.NewConnector(ariClient, reps.Bridge)

	return &Applications{
		Listener:        app.NewListener(reps, ariClient, app.NewHandler(ariClient, reps, connector)),
		Snoopy:          snoopy.New(),
		Conference:      app.NewConference(reps, ariClient),
		AsteriskAccount: app.NewAsteriskAccount(reps),
		User:            app.NewUser(reps),
		Connector:       connector,
	}
}
