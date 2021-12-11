package user

import (
	"encoding/json"

	"github.com/CyCoreSystems/ari/v5"
)

type User struct {
	Username        string   `redis:"username"      json:"username"`
	Email           string   `redis:"email"         json:"email"`
	AsteriskAccount string   `redis:"id"            json:"id"`
	ConferenceID    string   `redis:"conference_id" json:"conference_id"`
	NeedProtocol    bool     `redis:"need_protocol" json:"need_protocol"`
	SessionID       string   `redis:"session"`
	Channel         *ari.Key `redis:"channel"`
	Records         []string `redis:"records"`
	Texts           []string `redis:"texts"`
}

func (u *User) IsLess(then *User) bool {
	return u.SessionID < then.SessionID
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

type UserInfo struct {
	Name    string `json:"name"`
	ID      string `json:"id"`
	Channel string `json:"channel"`
}
