package api

import "time"

type User struct {
	Username string `json:"username" binding:"required"`
	Email    string `json:"email"`
	Path     string `json:"path" binding:"required"`
}

type TranslateRequest struct {
	StartTime time.Time `json:"start" binding:"required"`
	User      User      `json:"user" binding:"required"`
}

type SendProtocolRequest struct {
	ConferenceID string   `json:"conference_id" binding:"required"`
	To           []string `json:"to" binding:"required"`
}

