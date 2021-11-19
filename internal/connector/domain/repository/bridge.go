package repository

type Bridge interface {
	CreateBridge(hostUsername string, bridgeID string)
	GetForHost(hostUsername string) (string, error)
	DeleteBridge(bridgeID string) error
}
