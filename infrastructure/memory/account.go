package memory

import (
	"github.com/google/btree"
	"github.com/sirupsen/logrus"
	"protocall/domain/entity"
	"protocall/domain/repository"
	"sync"
)

type AsteriskAccountMemory struct {
	store *btree.BTree
	lock  *sync.RWMutex
}

func NewAsteriskAccount() *AsteriskAccountMemory {
	repo := &AsteriskAccountMemory{
		lock:  &sync.RWMutex{},
		store: btree.New(32),
	}

	repo.store.ReplaceOrInsert(&entity.AsteriskAccount{
		Username: "1233",
		Password: "technopark5535",
		UserID:   "",
	})
	repo.store.ReplaceOrInsert(&entity.AsteriskAccount{
		Username: "1234",
		Password: "technopark5535",
		UserID:   "",
	})
	repo.store.ReplaceOrInsert(&entity.AsteriskAccount{
		Username: "1235",
		Password: "technopark5535",
		UserID:   "",
	})
	return repo
}

func (a AsteriskAccountMemory) GetFree() *entity.AsteriskAccount {
	var freeAccount *entity.AsteriskAccount

	a.store.Ascend(func(item btree.Item) bool {
		account := item.(*entity.AsteriskAccount)
		if account == nil {
			return false
		}
		if account.UserID == "" {
			freeAccount = account
			return false
		}
		return true
	})

	return freeAccount
}

func (a AsteriskAccountMemory) Take(account string, userID string) {
	item := a.store.Get(&entity.AsteriskAccount{
		Username: account,
	})
	if item == nil {
		logrus.Error("Fail to take account")
	}

	accountItem := item.(*entity.AsteriskAccount)
	accountItem.UserID = userID
	a.store.ReplaceOrInsert(accountItem)
}

func (a AsteriskAccountMemory) Free(account string) {
	a.Take(account, "")
}

func (a AsteriskAccountMemory) Get(account string) *entity.AsteriskAccount {
	item := a.store.Get(&entity.AsteriskAccount{
		Username: account,
	})
	if item == nil {
		return nil
	}

	return item.(*entity.AsteriskAccount)
}

var _ repository.AsteriskAccountRepository = AsteriskAccountMemory{}
