package entity

import "github.com/google/btree"

type AsteriskAccount struct {
	Username string `json:"username"`
	Password string `json:"password"`
	UserID   string `json:"-"`
}

func (a AsteriskAccount) Less(then btree.Item) bool {
	return a.Username < then.(*AsteriskAccount).Username
}
