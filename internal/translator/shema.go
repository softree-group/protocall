package translator

import (
	"protocall/pkg/recognizer"
	"time"
)

type User struct {
	Username string `json:"username" binding:"required"`
	Record   string `json:"record" binding:"required"`
}

type TranslateRequest struct {
	Start time.Time `json:"start" binding:"required"`
	User  User      `json:"user" binding:"required"`
}

type TranslateRespone = recognizer.TextRespone
