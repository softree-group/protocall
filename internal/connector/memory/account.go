package memory

import (
	"protocall/internal/connector/config"
	"protocall/internal/connector/domain/entity"
	"protocall/internal/connector/domain/repository"
	"sync"

	"github.com/google/btree"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

type AsteriskAccountMemory struct {
	store   *btree.BTree
	userMap map[string]string
	lock    *sync.RWMutex
}

func NewAsteriskAccount() *AsteriskAccountMemory {
	repo := &AsteriskAccountMemory{
		lock:    &sync.RWMutex{},
		store:   btree.New(viper.GetInt(config.Participant)),
		userMap: map[string]string{},
	}
	return repo
}

func (a *AsteriskAccountMemory) GetFree() *entity.AsteriskAccount {
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

func (a *AsteriskAccountMemory) TakeAccount(account, userID string) {
	item := a.store.Get(&entity.AsteriskAccount{
		Username: account,
	})
	if item == nil {
		logrus.Error("Fail to take account")
	}

	a.userMap[account] = userID

	accountItem := item.(*entity.AsteriskAccount)
	accountItem.UserID = userID
	a.store.ReplaceOrInsert(accountItem)
}

func (a *AsteriskAccountMemory) FreeAccount(account string) {
	delete(a.userMap, account)
	a.TakeAccount(account, "")
}

func (a *AsteriskAccountMemory) GetAccount(account string) *entity.AsteriskAccount {
	item := a.store.Get(&entity.AsteriskAccount{
		Username: account,
	})
	if item == nil {
		return nil
	}

	return item.(*entity.AsteriskAccount)
}

func (a *AsteriskAccountMemory) SaveAccount(account entity.AsteriskAccount) {
	a.store.ReplaceOrInsert(&account)
}

func (a *AsteriskAccountMemory) Who(account string) string {
	return a.userMap[account]
}

var _ repository.AsteriskAccountRepository = &AsteriskAccountMemory{}
