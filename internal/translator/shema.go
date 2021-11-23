package translator

import (
	"protocall/pkg/yastt"
	"time"
)

type Record struct {
	Path   string        `json:"path"`
	Length time.Duration `json:"length"`
}

type User struct {
	Username    string    `json:"username"`
	ConnectTime time.Time `json:"join_time"`
	SessionID   string    `json:"session_id`
	Record      Record    `json:"record"`
	Text        string    `json:"text"`
}

type TranslateRequest struct {
	User `json:"user"`
}

type TranslateRespone = yastt.Chunk

type ConnectorRequest struct {
	SessionID string `json:"session_id"`
	Record    string `json:"record"`
	Text      string `json:"text"`
}
