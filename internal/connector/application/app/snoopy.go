package app

import (
	"context"

	"protocall/internal/connector/config"
	"protocall/internal/connector/domain/services"

	"github.com/CyCoreSystems/ari/v5"
	"github.com/CyCoreSystems/ari/v5/client/native"
	"github.com/CyCoreSystems/ari/v5/ext/record"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

type Snoopy struct {
	ari ari.Client
	bus services.Bus
}

func NewSnoopy(bus services.Bus) *Snoopy {
	ariClient, err := native.Connect(&native.Options{
		Application:  viper.GetString(config.ARISnoopyApplication),
		URL:          viper.GetString(config.ARIUrl),
		WebsocketURL: viper.GetString(config.ARIWebsocketURL),
		Username:     viper.GetString(config.ARIUser),
		Password:     viper.GetString(config.ARIPassword),
	})
	if err != nil {
		logrus.Fatal("Fail to connect snoopy app")
	}

	return &Snoopy{
		ari: ariClient,
		bus: bus,
	}
}

func (s *Snoopy) channelHandler(channel *ari.ChannelHandle, recordPath string, sessionID string) {
	sub := channel.Subscribe(ari.Events.All)
	//end := channel.Subscribe(ari.Events.StasisEnd)
	logrus.Debug("Record Path: ", recordPath)
	leave := s.bus.Subscribe("leave/" + sessionID)

	defer sub.Cancel()
	//defer end.Cancel()
	defer leave.Cancel()
	defer channel.Hangup()

	ctx := context.Background()
	rec := record.Record(ctx, channel)

	for {
		select {
		case event := <-sub.Events():
			logrus.Info("In SPY: ", event.GetType())
		case <-leave.Channel():
			logrus.Info("saving record for ", channel.ID())
			res := rec.Stop()

			err := res.Save(recordPath)
			if err != nil {
				logrus.Error("fail to save result record for channel ", channel.ID(), ". Error: ", err)
				return
			}
			logrus.Info("saved record for ", channel.ID())

			return
		}
	}
}

func (s *Snoopy) listen() {
	start := s.ari.Bus().Subscribe(nil, ari.Events.StasisStart)
	for event := range start.Events() {
		value := event.(*ari.StasisStart)

		channel := s.ari.Channel().Get(value.Key(ari.ChannelKey, value.Channel.ID))
		logrus.Info("snoop channel: ", channel.ID())

		go s.channelHandler(
			channel,
			value.Args[0],
			value.Args[1],
		)
	}
}

func (s *Snoopy) Snoop() {
	logrus.Info("Start snooping...")
	s.listen()
}
