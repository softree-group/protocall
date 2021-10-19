package entity

import (
	"github.com/google/btree"
)

type User struct {
	Username        string `json:"name"`
	Email           string `json:"-"`
	SessionID       string `json:"-"`
	AsteriskAccount string `json:"id"`
	Channel         string `json:"-"`
	ConferenceID    string `json:"conference_id"`
}

func (u User) Less(then btree.Item) bool {
	return u.SessionID < then.(*User).SessionID
}
