package memory

import (
	"protocall/domain/entity"
	"protocall/domain/repository"
	"sync"
)

type AsteriskAccountMemory struct {
	store []entity.AsteriskAccount
	lock *sync.RWMutex
}

func NewAsteriskAccountMemory() *AsteriskAccountMemory {
	return &AsteriskAccountMemory{
		lock: &sync.RWMutex{},
		store: []entity.AsteriskAccount{
			entity.AsteriskAccount{
				Username: "1235",
				Password: "technopark5535",
			},
			entity.AsteriskAccount{
				Username: "1233",
				Password: "technopark5535",
			},
			entity.AsteriskAccount{
				Username: "1234",
				Password: "technopark5535",
			},
		},
	}
}

func (a AsteriskAccountMemory) GetFree() *entity.AsteriskAccount {
	a.lock.RLock()
	defer a.lock.RUnlock()

	for _, account := range a.store {
		if account.UserID == "" {
			return &account
		}
	}
	return nil
}

func (a AsteriskAccountMemory) Take(account string, userID string) {
	for idx := 0; idx < len(a.store); idx++ {
		if a.store[idx].Username == account {
			a.lock.Lock()
			a.store[idx].UserID = userID
			a.lock.Unlock()
			return
		}
	}
}

func (a AsteriskAccountMemory) Free(account string) {
	for idx := 0; idx < len(a.store); idx++ {
		if a.store[idx].Username == account {
			a.lock.Lock()
			a.store[idx].UserID = ""
			a.lock.Unlock()
			return
		}
	}
}

var _ repository.AsteriskAccountRepository = AsteriskAccountMemory{}
