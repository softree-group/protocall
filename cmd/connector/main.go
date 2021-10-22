package main

import (
	"github.com/spf13/viper"
	"protocall/application"
	"protocall/config"
	"protocall/infrastructure"
	"protocall/interfaces/handlers"
)

func init() {
	config.InitConfig()
}

func main() {
	reps := infrastructure.New()
	apps := application.New(reps)

	apps.AsteriskAccount.Parse(viper.GetString(config.ARIAccountsFile))

	go handlers.ServeAPI(apps)
	go apps.Snoopy.Snoop()
	apps.Listener.Listen()
}
