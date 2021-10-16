package entity

import "github.com/google/btree"

type Conference struct {
	ID           string  `json:"id"`
	Participants []*User `json:"participants"`
	HostUserID   string  `json:"host_user_id"`
	BridgeID     string  `json:"-"`
}

func (c Conference) Less(then btree.Item) bool {
	return c.ID < then.(*Conference).ID
}
