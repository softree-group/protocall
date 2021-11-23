package entity

import "time"

type EventDefault struct {
	ConferenceID string
	Record       string
	Duration     time.Duration
	Text         string
	User         *User
}
