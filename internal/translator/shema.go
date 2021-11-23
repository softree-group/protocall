package translator

import (
	"protocall/pkg/yastt"
	"time"
)

type Record struct {
	Path   string        `json:"path" binding:"required"`
	Length time.Duration `json:"length" binding:"required"`
}

type User struct {
	Username    string    `json:"username" binding:"required"`
	ConnectTime time.Time `json:"join_time" binding:"required"`
	SessionID   string    `json:"session_id binding:"required"`
	Record      Record    `json:"record" binding:"required"`
	Text        string    `json:"text" binding:"required"`
}

type TranslateRequest struct {
	User `json:"user" binding:"required"`
}

type TranslateRespone = yastt.Chunk

type ConnectorRequest struct {
	SessionID string `json:"session_id" binding:"required"`
	Record    string `json:"record" binding:"required"`
	Text      string `json:"text" binding:"required"`
}
