package app

import (
	"encoding/json"
	"github.com/sirupsen/logrus"
	"io/ioutil"
	"os"
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

func (a AsteriskAccount) Parse(accountsFile string) {
	jsonFile, err := os.Open(accountsFile)
	if err != nil {
		logrus.Fatal("fail to open file ", accountsFile, ": ", err)
	}
	defer jsonFile.Close()

	bytes, _ := ioutil.ReadAll(jsonFile)

	var accounts entity.AsteriskAccounts

	err = json.Unmarshal(bytes, &accounts)
	if err != nil {
		logrus.Fatal("fail to parse file ", accountsFile, ": ", err)
	}

	for _, account := range accounts {
		a.reps.AsteriskAccount.Save(account)
	}
}

var _ applications.AsteriskAccount = AsteriskAccount{}
