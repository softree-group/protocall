package applications

import "protocall/internal/connector/domain/entity"

type Conference interface {
	StartConference(user *entity.User) (*entity.Conference, error)
	JoinToConference(user *entity.User, meetID string) (*entity.Conference, error)
	IsExist(meetID string) bool
	StartRecord(user *entity.User, meetID string) error
	UploadRecord(user *entity.User, meetID string) error
	TranslateRecord(user *entity.User, conference *entity.Conference) error
	CreateProtocol(conference *entity.Conference) error
	Get(meetID string) *entity.Conference
	StartRecordUser(user *entity.User, meetID string) error
	Delete(meetID string)
	RemoveParticipant(user *entity.User, meetID string)
}
