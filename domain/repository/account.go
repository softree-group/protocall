package repository

import "protocall/domain/entity"

type AsteriskAccountRepository interface {
	Get(account string) *entity.AsteriskAccount
	GetFree() *entity.AsteriskAccount
	Take(account string, userID string)
	Free(account string)
	Save(account entity.AsteriskAccount)
}
