package applications

import "protocall/domain/entity"

type Conference interface {
	StartConference(user *entity.User) (*entity.Conference, error)
	JoinToConference(user *entity.User, meetID string) (*entity.Conference, error)
	IsExist(meetID string) bool
}
