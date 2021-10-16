package memory

import (
	"github.com/google/btree"
	"protocall/domain/entity"
	"protocall/domain/repository"
)

type UserMemory struct {
	store *btree.BTree
}

func NewUser() *UserMemory {
	return &UserMemory{
		store: btree.New(32),
	}
}

func (u UserMemory) Find(sessionID string) *entity.User {
	item := u.store.Get(&entity.User{SessionID: sessionID})
	if item == nil {
		return nil
	}
	return item.(*entity.User)
}

func (u UserMemory) Save(user *entity.User) {
	u.store.ReplaceOrInsert(user)
}

func (u UserMemory) Delete(sessionID string) {
	u.store.Delete(&entity.User{SessionID: sessionID})
}

var _ repository.User = UserMemory{}
