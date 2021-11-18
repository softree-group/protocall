package app

import (
	"protocall/application/applications"
	"protocall/domain/entity"
	"protocall/domain/repository"

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
}

func NewListener(reps repository.Repositories,
	client ari.Client,
	handler applications.CallHandler,
	user applications.User,
	conference applications.Conference,
	account applications.AsteriskAccount,
	socket applications.Socket) *EventListener {
	return &EventListener{
		handler:    handler,
		ari:        client,
		reps:       reps,
		user:       user,
		conference: conference,
		account:    account,
		socket:     socket,
	}
}

func (e *EventListener) Listen() {
	logrus.Info("Start listen...")
	any := e.ari.Bus().Subscribe(nil, ari.Events.All)
	stasis := e.ari.Bus().Subscribe(nil, ari.Events.StasisStart)
	leftBridge := e.ari.Bus().Subscribe(nil, ari.Events.ChannelLeftBridge)

	for {
		select {
		case event := <-stasis.Events():
			data := event.(*ari.StasisStart)
			logrus.Info("new stasis", data)
			channel := e.ari.Channel().Get(ari.NewKey(ari.ChannelKey, data.Channel.ID))
			logrus.Info("Channel: ", channel.ID())
			channelData, err := channel.Data()
			if err != nil {
				logrus.Error("Fail to read data")
				channel.Hangup()
				continue
			}
			logrus.Info("Display name: ", channelData.Caller.Name)
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
			err = e.socket.PublishUserMessage(user, entity.SocketMessage{
				"event": "ready",
			})
			if err != nil {
				logrus.Error("fail to send socket message: ", err)
			}
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
