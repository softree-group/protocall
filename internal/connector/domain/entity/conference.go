package entity

import (
	"protocall/internal/connector/config"
	"time"

	"github.com/google/btree"
	"github.com/spf13/viper"
)

type Conference struct {
	ID           string       `json:"id"`
	Title        string       `json:"title"`
	Participants *btree.BTree `json:"participants"`
	HostUserID   string       `json:"host_user_id"`
	BridgeID     string       `json:"-"`
	IsRecording  bool         `json:"is_recording"`
	Start        time.Time
}

func NewConference(id, hostUser, title string) *Conference {
	return &Conference{
		ID:           id,
		Title:        title,
		Participants: btree.New(viper.GetInt(config.Participant)),
		HostUserID:   hostUser,
		BridgeID:     id,
		IsRecording:  false,
		Start:        time.Now(),
	}
}

func (c *Conference) Less(then btree.Item) bool {
	return c.ID < then.(*Conference).ID
}
