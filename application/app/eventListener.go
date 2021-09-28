package app

import (
	"github.com/CyCoreSystems/ari/v5"
	"github.com/sirupsen/logrus"
	"protocall/application/applications"
)

type EventListener struct {
	handler applications.CallHandler
	ari     ari.Client
}

func NewListener(client ari.Client, handler applications.CallHandler) *EventListener {
	return &EventListener{
		handler: handler,
		ari:     client,
	}
}

func (e EventListener) Listen() {
	logrus.Info("Start listen...")
	start := e.ari.Bus().Subscribe(nil, ari.Events.StasisStart)

	for {
		select {
		case event := <-start.Events():
			value := event.(*ari.StasisStart)

			channel := e.ari.Channel().Get(value.Key(ari.ChannelKey, value.Channel.ID))
			logrus.Info("catch channel: ", channel.ID())
			go e.handler.Handle(channel)
		}
	}

}

var _ applications.EventListener = &EventListener{}
