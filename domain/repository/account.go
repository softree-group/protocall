package repository

import "protocall/domain/entity"

type AsteriskAccountRepository interface {
	GetFree() *entity.AsteriskAccount
	Take(account string, userID string)
	Free(account string)
}
