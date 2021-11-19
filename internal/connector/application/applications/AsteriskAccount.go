package applications

import "protocall/domain/entity"

type AsteriskAccount interface {
	GetFree() *entity.AsteriskAccount
	Get(account string) *entity.AsteriskAccount
	Take(account string, userID string)
	Free(account string)
	Who(account string) string
}
