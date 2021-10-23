package repository

import (
	"protocall/domain/entity"
)

type User interface {
	FindUser(sessionID string) *entity.User
	SaveUser(user *entity.User)
	DeleteUser(sessionID string)
}
