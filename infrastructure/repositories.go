package infrastructure

import (
	"protocall/domain/repository"
	"protocall/infrastructure/memory"
)

func New() *repository.Repositories {
	return &repository.Repositories{
		Bridge:          &memory.Bridge{},
		AsteriskAccount: memory.NewAsteriskAccount(),
		User:            memory.NewUser(),
		Conference:      memory.NewConference(),
	}
}
