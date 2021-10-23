package main

import (
	"protocall/application"
	"protocall/infrastructure"
	"protocall/interfaces/handlers"
	"protocall/internal/config"
)

func main() {
	config.InitConfig()

	app := application.New(infrastructure.New())

	go handlers.ServeAPI(app)
	go app.Snoopy.Snoop()
	app.Listener.Listen()
}
