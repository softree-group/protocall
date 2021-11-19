package applications

import (
	"github.com/CyCoreSystems/ari/v5"
	"protocall/domain/entity"
)

type Connector interface {
	Connect(bridge *ari.BridgeHandle, channelID string) error
	CreateBridge(ID string) (*ari.BridgeHandle, error)
	CreateBridgeFrom(channel *ari.ChannelHandle) (*ari.BridgeHandle, error)
	Disconnect(bridgeID string, channel *ari.Key) error
	CallAndConnect(user *entity.User) (*ari.Key, error)
}
