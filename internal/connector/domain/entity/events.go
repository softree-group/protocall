package entity

import "time"

type EventDefault struct {
	ConferenceID string
	RecName string
	Duration time.Duration
	User *User
}
