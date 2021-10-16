package applications

import (
	"github.com/CyCoreSystems/ari/v5"
)

type Connector interface {
	Connect(bridge *ari.BridgeHandle, channelID string) error
	CreateBridge(ID string) (*ari.BridgeHandle, error)
	CreateBridgeFrom(channel *ari.ChannelHandle) (*ari.BridgeHandle, error)
	Disconnect(bridge *ari.BridgeHandle, channelID string) error
}
