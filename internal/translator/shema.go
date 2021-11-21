package translator

import (
	"protocall/pkg/recognizer"
	"time"
)

type User struct {
	Username    string    `json:"username" binding:"required"`
	ConnectTime time.Time `json:"join_time" binding:"required"`
	SessionID   string    `json:"session_id binding:"required"`
	Record      string    `json:"record" binding:"required"`
	Text        string    `json:"text" binding:"required"`
}

type TranslateRequest struct {
	User `json:"user" binding:"required"`
}

type TranslateRespone = recognizer.TextRespone
