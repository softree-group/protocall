package translator

import (
	"protocall/pkg/recognizer"
	"time"
)

type User struct {
	Username string `json:"username" binding:"required"`
	Path     string `json:"path" binding:"required"`
}

type TranslateRequest struct {
	ConferenceID string    `json:"conference_id" binding:"required"`
	StartTime    time.Time `json:"start" binding:"required"`
	User         User      `json:"user" binding:"required"`
}

type TranslateRespone = recognizer.TextRespone
