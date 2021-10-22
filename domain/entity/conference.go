package entity

import (
	"time"

	"github.com/google/btree"
)

type Conference struct {
	ID           string       `json:"id"`
	Participants *btree.BTree `json:"participants"`
	HostUserID   string       `json:"host_user_id"`
	BridgeID     string       `json:"-"`
	IsRecording  bool         `json:"is_recording"`
	Start        time.Time
}

func NewConference(id string, hostUser string) *Conference {
	return &Conference{
		ID:           id,
		Participants: btree.New(32),
		HostUserID:   hostUser,
		BridgeID:     id,
		IsRecording:  false,
		Start:        time.Now(),
	}
}

func (c Conference) Less(then btree.Item) bool {
	return c.ID < then.(*Conference).ID
}
