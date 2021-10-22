package app

import (
	"protocall/application/applications"
	"protocall/domain/repository"

	"github.com/CyCoreSystems/ari/v5"
	"github.com/sirupsen/logrus"
)

type EventListener struct {
	handler applications.CallHandler
	ari     ari.Client
	reps    repository.Repositories
}

func NewListener(reps repository.Repositories, client ari.Client, handler applications.CallHandler) *EventListener {
	return &EventListener{
		handler: handler,
		ari:     client,
		reps:    reps,
	}
}

func (e *EventListener) Listen() {
	logrus.Info("Start listen...")
	any := e.ari.Bus().Subscribe(nil, ari.Events.All)
	leftBridge := e.ari.Bus().Subscribe(nil, ari.Events.ChannelLeftBridge)

	for {
		select {
		case event := <-any.Events():
			logrus.Info("Event type: ", event.GetType())
		case event := <-leftBridge.Events():
			value := event.(*ari.ChannelLeftBridge)
			logrus.Info("Bridge ID: ", value.Bridge.ID)
			logrus.Info("LEFT: ", value.Bridge.ChannelIDs)
		}
	}

}

var _ applications.EventListener = &EventListener{}
