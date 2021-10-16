package app

import (
	"protocall/application/applications"
	"protocall/domain/entity"
	"protocall/domain/repository"
)

type AsteriskAccount struct {
	reps *repository.Repositories
}

func NewAsteriskAccount(reps *repository.Repositories) *AsteriskAccount {
	return &AsteriskAccount{
		reps: reps,
	}
}

func (a AsteriskAccount) GetFree() *entity.AsteriskAccount {
	return a.reps.AsteriskAccount.GetFree()
}

func (a AsteriskAccount) Get(account string) *entity.AsteriskAccount {
	return a.reps.AsteriskAccount.Get(account)
}

func (a AsteriskAccount) Take(account string, userID string) {
	a.reps.AsteriskAccount.Take(account, userID)
}

func (a AsteriskAccount) Free(account string) {
	a.reps.AsteriskAccount.Free(account)
}

var _ applications.AsteriskAccount = AsteriskAccount{}
