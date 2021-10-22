package memory

import "protocall/domain/repository"

type Bridge struct {
	bridgeID string
}

func NewBridge() *Bridge {
	return &Bridge{}
}

func (b *Bridge) CreateBridge(hostUsername string, bridgeID string) {
	b.bridgeID = bridgeID
}

func (b *Bridge) GetForHost(hostUsername string) (string, error) {
	return b.bridgeID, nil
}

func (b *Bridge) DeleteBridge(bridgeID string) error {
	return nil
}

var _ repository.Bridge = &Bridge{}
