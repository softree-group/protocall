package bridge

import (
	"context"
	"time"

	"protocall/internal/conference"
	"protocall/internal/translator"
	"protocall/internal/user"
	"protocall/pkg/logger"

	"github.com/CyCoreSystems/ari/v5"
	"github.com/CyCoreSystems/ari/v5/ext/record"
)

func (a *Application) channelHandler(channel *ari.ChannelHandle, recordPath, sessionID string) {
	defer channel.Hangup()

	sub := channel.Subscribe(ari.Events.All)
	defer sub.Cancel()

	leave := a.bus.Subscribe("leave/" + sessionID)
	defer leave.Cancel()

	rec := record.Record(context.Background(), channel)
	a.bus.Publish(
		"startRecord",
		conference.Event{
			Record: &translator.Record{
				Path: recordPath,
			},
			User: &user.User{
				SessionID: sessionID,
			},
		})
	startedTime := time.Now()

	for {
		select {
		case event := <-sub.Events():
			logger.L.Info("In SPY: ", event.GetType())
		case <-leave.Channel():
			logger.L.Info("saving record for ", channel.ID())
			res := rec.Stop()

			err := res.Save(recordPath)
			if err != nil {
				logger.L.Error("fail to save result record for channel ", channel.ID(), ". Error: ", err)
				return
			}

			a.bus.Publish("saved", conference.Event{
				Record: &translator.Record{
					Path:   recordPath,
					Length: time.Since(startedTime),
				},
				User: &user.User{
					SessionID: sessionID,
				},
			})
			logger.L.Info("saved record for ", channel.ID())
			return
		}
	}
}

func (a *Application) Snoopy() {
	start := a.ari.Bus().Subscribe(nil, ari.Events.StasisStart)
	for event := range start.Events() {
		value := event.(*ari.StasisStart)

		channel := a.ari.Channel().Get(value.Key(ari.ChannelKey, value.Channel.ID))
		logger.L.Info("snoop channel: ", channel.ID())

		go a.channelHandler(
			channel,
			value.Args[0],
			value.Args[1],
		)
	}
}
