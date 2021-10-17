package repository

import (
	"protocall/domain/entity"
)

type User interface {
	Find(sessionID string) *entity.User
	Save(user *entity.User)
	Delete(sessionID string)
}
