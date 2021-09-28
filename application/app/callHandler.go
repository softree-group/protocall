package app

import (
	"github.com/CyCoreSystems/ari/v5"
	"github.com/sirupsen/logrus"
	"protocall/application/applications"
	"protocall/domain/repository"
)

type CallHandler struct {
	reps *repository.Repositories
	ari  ari.Client
}

func NewHandler(client ari.Client, reps *repository.Repositories) *CallHandler {
	return &CallHandler{
		reps: reps,
		ari:  client,
	}
}

func (c CallHandler) Handle(channel *ari.ChannelHandle) {
	bridgeID, err := c.reps.Bridge.GetForHost("some")
	if err != nil {
		logrus.Error("no bridge")
		return
	}

	bridge := c.ari.Bridge().Get(channel.Key().New(ari.BridgeKey, bridgeID))
	if bridge == nil {
		logrus.Error("no bridge ", bridgeID)
		return
	}

	_ = channel.Answer()

	//ctx := context.Background()

	err = bridge.AddChannel(channel.ID())
	if err != nil {
		logrus.Error("cannot connect ", channel.ID(), " to bridge ", bridgeID)
		return
	}

	end := channel.Subscribe(ari.Events.ChannelLeftBridge)
	anyEvent := channel.Subscribe(ari.Events.All)

	_, err = channel.Snoop("snoop_"+channel.ID(), &ari.SnoopOptions{
		App:     "snoopy",
		AppArgs: channel.ID(),
		Spy:     "in",
		Whisper: "both",
	})

	if err != nil {
		logrus.Error("Fail to snoop: ", err)
	}

	defer end.Cancel()
	defer anyEvent.Cancel()

	for {
		select {
		case e := <-anyEvent.Events():
			logrus.Info("EVENT TYPE: ", e.GetType(), " for ", channel.ID())
		case <-end.Events():
			logrus.Info("channel ", channel.ID(), " hangup")
			return
		}
	}
}

var _ applications.CallHandler = &CallHandler{}
