package app

import (
	"protocall/application/applications"
	"protocall/domain/entity"
	"protocall/domain/repository"
)

type User struct {
	reps *repository.Repositories
}

func (u User) Find(sessionID string) *entity.User {
	return u.reps.User.Find(sessionID)
}

func (u User) Save(user *entity.User) {
	u.reps.User.Save(user)
}

func (u User) Delete(sessionID string) {
	u.reps.User.Delete(sessionID)
}

func NewUser(reps *repository.Repositories) *User {
	return &User{reps: reps}
}

var _ applications.User = &User{}
