package translator

import (
	"protocall/pkg/recognizer"
	"time"
)

type TranslateRequest struct {
	StartTime time.Time `json:"start" binding:"required"`
	User      struct {
		Username string `json:"username" binding:"required"`
		Email    string `json:"email"`
		Path     string `json:"path" binding:"required"`
	} `json:"user" binding:"required"`
}

type TranslateRespone = recognizer.TextRespone
