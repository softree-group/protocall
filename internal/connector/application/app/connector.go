package app

import (
	"fmt"

	"protocall/internal/connector/application/applications"
	"protocall/internal/connector/config"
	"protocall/internal/connector/domain/entity"
	"protocall/internal/connector/domain/repository"

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

func (c *Connector) getBridge(id string) *ari.BridgeHandle {
	key := &ari.Key{
		Kind:                 ari.BridgeKey,
		ID:                   id,
		Node:                 "",
		Dialog:               "",
		App:                  viper.GetString(config.ARIApplication),
		XXX_NoUnkeyedLiteral: struct{}{},
		XXX_unrecognized:     nil,
		XXX_sizecache:        0,
	}

	return c.ari.Bridge().Get(key)
}

func (c *Connector) CreateBridge(id string) (*ari.BridgeHandle, error) {
	key := &ari.Key{
		Kind:                 ari.BridgeKey,
		ID:                   id,
		Node:                 "",
		Dialog:               "",
		App:                  viper.GetString(config.ARIApplication),
		XXX_NoUnkeyedLiteral: struct{}{},
		XXX_unrecognized:     nil,
		XXX_sizecache:        0,
	}

	return c.ari.Bridge().Create(key, "video_sfu", key.ID)
}

func (c *Connector) CallAndConnect(user *entity.User) (*ari.Key, error) {
	bridgeID := user.ConferenceID
	account := user.AsteriskAccount
	if user.Channel != nil {
		ch := c.ari.Channel().Get(user.Channel)
		if ch != nil {
			err := ch.Hangup()
			if err != nil {
				logrus.Error("fail to hangup: ", err)
			}
		}
	}
	bridge := c.getBridge(user.ConferenceID)
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
	err := c.ari.Channel().Get(channel).Hangup()
	if err != nil {
		logrus.Error("fail to delete channel: ", err)
	}
	return err
}

var _ applications.Connector = &Connector{}
