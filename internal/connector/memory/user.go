package memory

import (
	"protocall/internal/connector/config"
	"protocall/internal/connector/domain/entity"
	"protocall/internal/connector/domain/repository"

	"github.com/google/btree"
	"github.com/spf13/viper"
)

type UserMemory struct {
	store *btree.BTree
}

func NewUser() *UserMemory {
	return &UserMemory{
		store: btree.New(viper.GetInt(config.Participant)),
	}
}

func (u UserMemory) FindUser(sessionID string) *entity.User {
	item := u.store.Get(&entity.User{SessionID: sessionID})
	if item == nil {
		return nil
	}
	return item.(*entity.User)
}

func (u UserMemory) SaveUser(user *entity.User) {
	u.store.ReplaceOrInsert(user)
}

func (u UserMemory) DeleteUser(sessionID string) {
	u.store.Delete(&entity.User{SessionID: sessionID})
}

var _ repository.User = UserMemory{}
