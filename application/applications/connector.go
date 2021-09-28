package applications

import "github.com/CyCoreSystems/ari/v5"

type Connector interface {
	Connect(bridge *ari.BridgeHandle, channelID string) error
	CreateBridge() (*ari.BridgeHandle, error)
	Disconnect(bridge *ari.BridgeHandle, channelID string) error
}
