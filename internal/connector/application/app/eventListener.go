package app

import (
	"protocall/internal/connector/application/applications"
	"protocall/internal/connector/domain/entity"
	"protocall/internal/connector/domain/repository"
	"protocall/internal/connector/domain/services"

	"github.com/CyCoreSystems/ari/v5"
	"github.com/sirupsen/logrus"
)

type EventListener struct {
	handler    applications.CallHandler
	ari        ari.Client
	reps       repository.Repositories
	user       applications.User
	conference applications.Conference
	account    applications.AsteriskAccount
	socket     applications.Socket
	bus        services.Bus
}

func NewListener(reps repository.Repositories,
	client ari.Client,
	handler applications.CallHandler,
	user applications.User,
	conference applications.Conference,
	account applications.AsteriskAccount,
	socket applications.Socket,
	bus services.Bus) *EventListener {
	return &EventListener{
		handler:    handler,
		ari:        client,
		reps:       reps,
		user:       user,
		conference: conference,
		account:    account,
		socket:     socket,
		bus:        bus,
	}
}

func (e *EventListener) Listen() {
	logrus.Info("Start listen...")
	stasis := e.ari.Bus().Subscribe(nil, ari.Events.StasisStart)
	leftBridge := e.ari.Bus().Subscribe(nil, ari.Events.ChannelLeftBridge)

	for {
		select {
		case event := <-stasis.Events():
			data := event.(*ari.StasisStart)
			channel := e.ari.Channel().Get(ari.NewKey(ari.ChannelKey, data.Channel.ID))
			channelData, err := channel.Data()
			if err != nil {
				logrus.Error("Fail to read data")
				channel.Hangup()
				continue
			}
			sessionID := e.account.Who(channelData.Caller.Number)
			if sessionID == "" {
				logrus.Warn("Free account ", channelData.Caller.Number)
				channel.Hangup()
				continue
			}
			user := e.user.Find(sessionID)
			if user == nil {
				logrus.Warn("No user")
				e.account.Free(channelData.Caller.Number)
				e.user.Delete(sessionID)
				channel.Hangup()
				continue
			}

			if user.Channel != nil {
				e.bus.Publish("leave/"+user.SessionID, "")
			}

			user.Channel = channel.Key()
			e.user.Save(user)

			conference := e.conference.Get(user.ConferenceID)
			if conference == nil {
				logrus.Warn("no conference: ", user.ConferenceID)
				e.socket.PublishUserMessage(user, entity.SocketMessage{
					"event":   "fail",
					"message": "no conference",
				})
				channel.Hangup()
				continue
			}

			if conference.IsRecording {
				err = e.conference.StartRecordUser(user, user.ConferenceID)
				if err != nil {
					logrus.Error("fail to start record: ", err)
				}
			}

			err = channel.Continue("conf", user.ConferenceID, 0)
			if err != nil {
				logrus.Info("Fail to continue: ", err)
				e.socket.PublishUserMessage(user, entity.SocketMessage{
					"event":   "fail",
					"message": "fail to continue",
				})
				channel.Hangup()
				continue
			}

			err = e.socket.PublishConnectedEvent(user)
			if err != nil {
				logrus.Error("fail to send socket message for connected event: ", err)
			}
			err = e.socket.PublishUserMessage(user, entity.SocketMessage{
				"event": "ready",
			})
			if err != nil {
				logrus.Error("fail to send socket message for ready event: ", err)
			}
		case event := <-leftBridge.Events():
			value := event.(*ari.ChannelLeftBridge)
			logrus.Info("Bridge ID: ", value.Bridge.ID)
			logrus.Info("LEFT: ", value.Bridge.ChannelIDs)
		}
	}
}

var _ applications.EventListener = &EventListener{}
