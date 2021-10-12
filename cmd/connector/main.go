package main

import (
	"protocall/application"
	"protocall/config"
	"protocall/infrastructure"
)

func init() {
	config.InitConfig()
}

func main() {
	reps := infrastructure.New()
	apps := application.New(reps)

	go apps.Snoopy.Snoop()
	apps.Listener.Listen()
}
