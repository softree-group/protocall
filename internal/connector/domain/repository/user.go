package repository

import (
	"protocall/internal/connector/domain/entity"
)

type User interface {
	FindUser(sessionID string) *entity.User
	SaveUser(user *entity.User)
	DeleteUser(sessionID string)
}
