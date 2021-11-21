package main

import (
	"protocall/internal/connector/application"
	"protocall/internal/connector/config"
	"protocall/internal/connector/handlers"
)

func main() {
	config.InitConfig()

	app := application.New(handlers.NewHandler())

	go handlers.ServeAPI(app)
	go app.Snoopy.Snoop()
	go app.ApplicationEventListener.Listen()
	app.Listener.Listen()
}
