package app

import (
	"protocall/application/applications"
	"protocall/domain/entity"
	"protocall/domain/repository"
)

type User struct {
	reps repository.User
}

func NewUser(reps repository.User) *User {
	return &User{reps: reps}
}

func (u *User) Find(sessionID string) *entity.User {
	return u.reps.FindUser(sessionID)
}

func (u *User) Save(user *entity.User) {
	u.reps.SaveUser(user)
}

func (u *User) Delete(sessionID string) {
	u.reps.DeleteUser(sessionID)
}

var _ applications.User = &User{}
