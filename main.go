package main

import (
	"github.com/sirupsen/logrus"
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

	bridge, err := apps.Connector.CreateBridge()
	if err != nil {
		logrus.Fatal("fail to create bridge ", err)
	}

	reps.Bridge.Create("some", bridge.ID())

	go apps.Snoopy.Snoop()
	apps.Listener.Listen()
}
