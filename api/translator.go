package api

import "time"

type User struct {
	Username string `json:"username" binding:"required"`
	Email    string `json:"email"`
	Path     string `json:"path" binding:"required"`
}

type TranslatorRequest struct {
	ConfID    string    `json:"conf_id" binding:"required"`
	StartTime time.Time `json:"start" binding:"required"`
	User      User
}
