package bridge

type Storage struct {
	bridgeID string
}

func NewStorage() *Storage {
	return &Storage{}
}

func (b *Storage) CreateBridge(hostUsername, bridgeID string) {
	b.bridgeID = bridgeID
}

func (b *Storage) GetForHost(hostUsername string) (string, error) {
	return b.bridgeID, nil
}

func (b *Storage) DeleteBridge(bridgeID string) error {
	return nil
}
