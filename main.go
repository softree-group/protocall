package main

import (
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

	go handlers.ServeAPI(apps)
	go apps.Snoopy.Snoop()
	apps.Listener.Listen()
}
