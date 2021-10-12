package app

import (
	"github.com/CyCoreSystems/ari/v5"
	"github.com/sirupsen/logrus"
	"protocall/application/applications"
	"protocall/domain/repository"
)

type CallHandler struct {
	reps      *repository.Repositories
	ari       ari.Client
	connector applications.Connector
}

func NewHandler(client ari.Client, reps *repository.Repositories, connector applications.Connector) *CallHandler {
	return &CallHandler{
		reps:      reps,
		ari:       client,
		connector: connector,
	}
}

func (c CallHandler) Handle(channel *ari.ChannelHandle) {
	bridgeID, err := c.reps.Bridge.GetForHost("some")
	logrus.Info("GET: ", bridgeID)
	if err != nil {
		logrus.Error("no bridge")
		return
	}
	var bridge *ari.BridgeHandle
	if bridgeID == "" {
		bridge, err = c.connector.CreateBridgeFrom(channel)
		logrus.Info("CREATED: ", bridge.ID())
		if err != nil {
			logrus.Error("IN CREATED: ", err)
		}
	} else {
		bridge = c.ari.Bridge().Get(channel.Key().New(ari.BridgeKey, bridgeID))
		logrus.Info("GET: ", bridge.ID())
	}

	_ = channel.Answer()

	//ctx := context.Background()

	err = bridge.AddChannel(channel.ID())
	if err != nil {
		logrus.Error("cannot connect ", channel.ID(), " to bridge ", bridge.ID(), " err: ", err)
		return
	}

	end := channel.Subscribe(ari.Events.ChannelLeftBridge)
	anyEvent := channel.Subscribe(ari.Events.All)

	data, err := channel.Data()
	logrus.Infof("DATA: %v+", data)

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
