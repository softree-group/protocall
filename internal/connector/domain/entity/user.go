package entity

import (
	"encoding/json"
	"github.com/CyCoreSystems/ari/v5"
	"github.com/google/btree"
)

type User struct {
	Username        string   `json:"name"`
	Email           string   `json:"email"`
	SessionID       string   `json:"-"`
	AsteriskAccount string   `json:"id"`
	Channel         *ari.Key `json:"-"`
	ConferenceID    string   `json:"conference_id"`
	RecordPath      string
	NeedProtocol    bool `json:"need_protocol"`
}

func (u *User) Less(then btree.Item) bool {
	return u.SessionID < then.(*User).SessionID
}

func (u User) MarshalJSON() ([]byte, error) {
	channel := ""
	if u.Channel != nil {
		channel = u.Channel.ID
	}
	return json.Marshal(&map[string]interface{}{
		"name":    u.Username,
		"email":   u.Email,
		"id":      u.AsteriskAccount,
		"channel": channel,
	})
}
