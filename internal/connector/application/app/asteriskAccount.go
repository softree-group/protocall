package app

import (
	"encoding/json"
	"io"
	"os"

	"protocall/internal/connector/application/applications"
	"protocall/internal/connector/domain/entity"
	"protocall/internal/connector/domain/repository"

	"github.com/sirupsen/logrus"
)

type AsteriskAccount struct {
	reps repository.AsteriskAccountRepository
}

func (a *AsteriskAccount) parse(accountsFile string) {
	jsonFile, err := os.Open(accountsFile)
	if err != nil {
		logrus.Fatal("fail to open file ", accountsFile, ": ", err)
	}
	defer jsonFile.Close()

	bytes, _ := io.ReadAll(jsonFile)

	var accounts entity.AsteriskAccounts

	err = json.Unmarshal(bytes, &accounts)
	if err != nil {
		logrus.Fatal("fail to parse file ", accountsFile, ": ", err)
	}

	for _, account := range accounts {
		a.reps.SaveAccount(account)
	}
}

func NewAsteriskAccount(reps repository.AsteriskAccountRepository, accountsFile string) *AsteriskAccount {
	r := &AsteriskAccount{
		reps: reps,
	}
	r.parse(accountsFile)
	return r
}

func (a *AsteriskAccount) GetFree() *entity.AsteriskAccount {
	return a.reps.GetFree()
}

func (a *AsteriskAccount) Get(account string) *entity.AsteriskAccount {
	return a.reps.GetAccount(account)
}

func (a *AsteriskAccount) Take(account, userID string) {
	a.reps.TakeAccount(account, userID)
}

func (a *AsteriskAccount) Free(account string) {
	a.reps.FreeAccount(account)
}

func (a *AsteriskAccount) Who(account string) string {
	return a.reps.Who(account)
}

var _ applications.AsteriskAccount = &AsteriskAccount{}
