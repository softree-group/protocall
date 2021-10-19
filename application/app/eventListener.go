package app

import (
	"github.com/CyCoreSystems/ari/v5"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"protocall/application/applications"
	"protocall/config"
	"protocall/domain/repository"
)

type EventListener struct {
	handler applications.CallHandler
	ari     ari.Client
	reps    *repository.Repositories
}

func NewListener(reps *repository.Repositories, client ari.Client, handler applications.CallHandler) *EventListener {
	return &EventListener{
		handler: handler,
		ari:     client,
		reps:    reps,
	}
}

func (e EventListener) Listen() {
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
			if len(value.Bridge.ChannelIDs) == 0 {
				err := e.ari.Bridge().Delete(&ari.Key{
					Kind: ari.BridgeKey,
					ID:   value.Bridge.ID,
					App:  viper.GetString(config.ARIApplication),
				})
				if err != nil {
					logrus.Error("fail to delete bridge: ", err)
				}

				conference := e.reps.Conference.Get(value.Bridge.ID)
				for _, user := range conference.Participants {
					e.reps.AsteriskAccount.Free(user.AsteriskAccount)
					e.reps.User.Delete(user.SessionID)
				}

				e.reps.Conference.Delete(value.Bridge.ID)
			}
		}
	}

}

var _ applications.EventListener = &EventListener{}
