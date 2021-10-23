package repository

import "protocall/domain/entity"

type AsteriskAccountRepository interface {
	GetAccount(account string) *entity.AsteriskAccount
	GetFree() *entity.AsteriskAccount
	TakeAccount(account string, userID string)
	FreeAccount(account string)
	SaveAccount(account entity.AsteriskAccount)
}
