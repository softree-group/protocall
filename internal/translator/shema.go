package translator

import (
	"time"

	"protocall/pkg/yastt"
)

type Record struct {
	URI    string        `json:"uri"`
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
	Record    Record `json:"record"`
	Text      string `json:"text"`
}
