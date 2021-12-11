package bridge

import (
	"protocall/internal/socket"
	"protocall/pkg/logger"

	"github.com/CyCoreSystems/ari/v5"
)

func (a *Application) Upstream() {
	logger.L.Info("Start listen...")
	stasis := a.ari.Bus().Subscribe(nil, ari.Events.StasisStart)
	leftBridge := a.ari.Bus().Subscribe(nil, ari.Events.ChannelLeftBridge)

	for {
		select {
		case event := <-stasis.Events():
			data := event.(*ari.StasisStart)
			channel := a.ari.Channel().Get(ari.NewKey(ari.ChannelKey, data.Channel.ID))
			channelData, err := channel.Data()
			if err != nil {
				logger.L.Error("Fail to read data")
				channel.Hangup()
				continue
			}
			sessionID := a.account.Who(channelData.Caller.Number)
			if sessionID == "" {
				logger.L.Warn("Free account ", channelData.Caller.Number)
				channel.Hangup()
				continue
			}
			user := a.user.Find(sessionID)
			if user == nil {
				logger.L.Warn("No user")
				a.account.Free(channelData.Caller.Number)
				a.user.Delete(sessionID)
				channel.Hangup()
				continue
			}

			if user.Channel != nil {
				a.bus.Publish("leave/"+user.SessionID, "")
			}

			user.Channel = channel.Key()
			a.user.Save(user)

			conference := a.conference.Get(user.ConferenceID)
			if conference == nil {
				logger.L.Warn("no conference: ", user.ConferenceID)
				a.socket.PublishUserMessage(user, socket.Message{
					"event":   "fail",
					"message": "no conference",
				})
				channel.Hangup()
				continue
			}

			if conference.IsRecording {
				err = a.conference.StartRecordUser(user, user.ConferenceID)
				if err != nil {
					logger.L.Error("fail to start record: ", err)
				}
			}

			err = channel.Continue("conf", user.ConferenceID, 0)
			if err != nil {
				logger.L.Info("Fail to continue: ", err)
				a.socket.PublishUserMessage(user, socket.Message{
					"event":   "fail",
					"message": "fail to continue",
				})
				channel.Hangup()
				continue
			}

			err = a.socket.PublishConnectedEvent(user)
			if err != nil {
				logger.L.Error("fail to send socket message for connected event: ", err)
			}
			err = a.socket.PublishUserMessage(user, socket.Message{
				"event": "ready",
			})
			if err != nil {
				logger.L.Error("fail to send socket message for ready event: ", err)
			}
		case event := <-leftBridge.Events():
			value := event.(*ari.ChannelLeftBridge)
			logger.L.Info("Bridge ID: ", value.Bridge.ID)
			logger.L.Info("LEFT: ", value.Bridge.ChannelIDs)
		}
	}
}
