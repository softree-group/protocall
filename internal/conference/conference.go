package conference

import (
	"time"

	"protocall/internal/user"
)

type Conference struct {
	ID           string      `json:"id"`
	Participants []user.User `json:"participants"`
	HostUserID   string      `json:"host_user_id"`
	IsRecording  bool        `json:"is_recording"`
	Start        time.Time
	BridgeID     string
}

func NewConference(id, hostUser string) *Conference {
	return &Conference{
		ID:          id,
		HostUserID:  hostUser,
		BridgeID:    id,
		IsRecording: false,
		Start:       time.Now(),
	}
}

func (c *Conference) Less(then Conference) bool {
	return c.ID < then.ID
}

type ConferenceInfo struct {
	ID           string          `json:"id"`
	HostID       string          `json:"host_id"`
	Participants []user.UserInfo `json:"participants"`
	IsRecording  bool            `json:"is_recording"`
	StartedAt    int64           `json:"started_at"`
}
