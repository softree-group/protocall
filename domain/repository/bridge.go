package repository

type Bridge interface {
	Create(hostUsername string, bridgeID string)
	GetForHost(hostUsername string) (string, error)
	Delete(bridgeID string) error
}
