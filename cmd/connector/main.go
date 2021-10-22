package main

import (
	"protocall/application"
	"protocall/infrastructure"
	"protocall/interfaces/handlers"
	"protocall/internal/config"
)

func init() {
	config.InitConfig()
}

func main() {
	app := application.New(infrastructure.New())

	go handlers.ServeAPI(app)
	go app.Snoopy.Snoop()
	app.Listener.Listen()
}
