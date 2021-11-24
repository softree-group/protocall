package entity

import "time"

type Record struct {
	Path   string
	URI    string
	Length time.Duration
}

type EventDefault struct {
	Record       *Record
	User         *User
	ConferenceID string
	Text         string
}
