package applications

import "protocall/internal/connector/domain/entity"

type AsteriskAccount interface {
	GetFree() *entity.AsteriskAccount
	Get(account string) *entity.AsteriskAccount
	Take(account string, userID string)
	Free(account string)
	Who(account string) string
}
