package app

import (
	"github.com/CyCoreSystems/ari/v5"
	"protocall/application/applications"
)

type Connector struct {
	ari ari.Client
}

func NewConnector(client ari.Client) *Connector {
	return &Connector{ari: client}
}

func (c Connector) CreateBridge() (*ari.BridgeHandle, error) {

	key := &ari.Key{
		Kind:                 ari.BridgeKey,
		ID:                   "some",
		Node:                 "",
		Dialog:               "",
		App:                  "protocall",
		XXX_NoUnkeyedLiteral: struct{}{},
		XXX_unrecognized:     nil,
		XXX_sizecache:        0,
	}

	bridge := c.ari.Bridge().Get(key)
	var err error
	if bridge == nil {
		bridge, err = c.ari.Bridge().Create(key, "mixing", key.ID)
	}

	if err != nil {
		return nil, err
	}

	return bridge, nil
}

func (c Connector) Connect(bridge *ari.BridgeHandle, channelID string) error {
	return bridge.AddChannel(channelID)
}

func (c Connector) Disconnect(bridge *ari.BridgeHandle, channelID string) error {
	return bridge.RemoveChannel(channelID)
}

var _ applications.Connector = &Connector{}
