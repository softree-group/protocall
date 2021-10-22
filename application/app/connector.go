package app

import (
	"errors"
	"fmt"
	"protocall/application/applications"
	"protocall/domain/repository"
	"protocall/internal/config"

	"github.com/CyCoreSystems/ari/v5"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

type Connector struct {
	ari         ari.Client
	bridgeStore repository.Bridge
}

func NewConnector(client ari.Client, bridgeStore repository.Bridge) *Connector {
	return &Connector{ari: client, bridgeStore: bridgeStore}
}

func (c *Connector) CreateBridgeFrom(channel *ari.ChannelHandle) (*ari.BridgeHandle, error) {

	key := channel.Key().New(ari.BridgeKey, channel.ID())

	bridge, err := c.ari.Bridge().Create(key, "video_sfu", key.ID)

	if err != nil {
		return nil, err
	}

	c.bridgeStore.CreateBridge(channel.ID(), bridge.ID())

	return bridge, nil
}

func (c *Connector) HasBridge() bool {
	bID, _ := c.bridgeStore.GetForHost("some")
	return bID != ""
}

func (c *Connector) getBridge(ID string) *ari.BridgeHandle {
	key := &ari.Key{
		Kind:                 ari.BridgeKey,
		ID:                   ID,
		Node:                 "",
		Dialog:               "",
		App:                  viper.GetString(config.ARIApplication),
		XXX_NoUnkeyedLiteral: struct{}{},
		XXX_unrecognized:     nil,
		XXX_sizecache:        0,
	}

	return c.ari.Bridge().Get(key)
}

func (c *Connector) CreateBridge(ID string) (*ari.BridgeHandle, error) {
	key := &ari.Key{
		Kind:                 ari.BridgeKey,
		ID:                   ID,
		Node:                 "",
		Dialog:               "",
		App:                  viper.GetString(config.ARIApplication),
		XXX_NoUnkeyedLiteral: struct{}{},
		XXX_unrecognized:     nil,
		XXX_sizecache:        0,
	}

	return c.ari.Bridge().Create(key, "mixing", key.ID)
}

func (c *Connector) CallAndConnect(account, bridgeID string) (*ari.Key, error) {
	bridge := c.getBridge(bridgeID)
	if bridge == nil {
		return nil, fmt.Errorf("bridge %s does not exist", bridgeID)
	}

	clientChannel, err := c.createCallInternal(account)
	if err != nil {
		return nil, err
	}

	err = c.waitUp(clientChannel)
	if err != nil {
		return nil, err
	}

	err = bridge.AddChannel(clientChannel.ID())
	if err != nil {
		return nil, err
	}

	return clientChannel.Key(), nil
}

func (c *Connector) Connect(bridge *ari.BridgeHandle, channelID string) error {
	return bridge.AddChannel(channelID)
}

func (c *Connector) Disconnect(bridgeID string, channel *ari.Key) error {
	bridge := c.getBridge(bridgeID)
	if bridge == nil {
		return errors.New("no bridge")
	}
	err := bridge.RemoveChannel(channel.ID)
	if err != nil {
		logrus.Error("fail to remove from channel: ", err)
	}

	err = c.ari.Channel().Get(channel).Hangup()
	if err != nil {
		logrus.Error("fail to delete channel: ", err)
	}
	return err
}

var _ applications.Connector = &Connector{}
