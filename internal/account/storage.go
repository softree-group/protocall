package account

import (
	"errors"
	"fmt"
	"protocall/pkg/logger"

	"github.com/gomodule/redigo/redis"
)

type Storage struct {
	*redis.Pool
}

func NewStorage(conn *redis.Pool) *Storage {
	return &Storage{conn}
}

const pattern = "account:*"

func (s *Storage) GetFree() (*Account, error) {
	var freeAccount *Account

	for i := 0; ; {
		arr, err := redis.Values(s.Do("SCAN", i, "MATCH"))
		if err != nil {
			return nil, fmt.Errorf("error retrieving '%s' keys", pattern)
		}
		if i == 0 {
			break
		}
	}

	return freeAccount, nil
}

func (s *Storage) TakeAccount(account, userID string) {
	item := a.store.Get(&Account{
		Username: account,
	})
	if item == nil {
		logger.L.Error("Fail to take account")
	}

	a.userMap[account] = userID

	accountItem := item.(*Account)
	accountItem.UserID = userID
	a.store.ReplaceOrInsert(accountItem)
}

// FreeAccount
func (s *Storage) FreeAccount(id string) error {
	account, err := a.GetAccount(id)
	if err != nil {
		return "", err
	}

	account.UserID = ""

	if _, err := a.Do("DEL", "account:"+id); err != nil {
		return err
	}
	return nil
}

var errNoUser = errors.New("cannot find user by id")

// GetAccount returns account by ID.
func (s *Storage) GetAccount(id string) (*Account, error) {
	values, err := redis.Values(a.Do("HGETALL", "account:"+id))
	if err != nil {
		return nil, err
	} else if len(values) == 0 {
		return nil, errNoUser
	}

	var account Account
	if err := redis.ScanStruct(values, &account); err != nil {
		return nil, err
	}

	return &account, nil
}

// SaveAccount add
func (s *Storage) SaveAccount(id string, account Account) error {
	if _, err := a.Do(
		"HMSET",
		"account:"+id,
		"id",
		account.UserID,
		"username",
		account.Username,
		"password",
		account.Password,
	); err != nil {
		return err
	}
	return nil
}

// Who returns user ID by account ID.
func (a *Account) Who(id string) (string, error) {
	account, err := a.GetAccount(id)
	if err != nil {
		return "", err
	}
	return account.UserID, nil
}
