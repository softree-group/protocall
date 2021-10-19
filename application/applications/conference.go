package applications

import (
	"protocall/domain/entity"
)

type Conference interface {
	StartConference(user *entity.User) (*entity.Conference, error)
	JoinToConference(user *entity.User, meetID string) (*entity.Conference, error)
	IsExist(meetID string) bool
	StartRecord(user *entity.User, meetID string) error
	Get(meetID string) *entity.Conference
	StartRecordUser(user *entity.User, meetID string) error
}
