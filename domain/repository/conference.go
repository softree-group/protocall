package repository

import "protocall/domain/entity"

type Conference interface {
	GetConference(conferenceID string) *entity.Conference
	SaveConference(conference *entity.Conference)
	DeleteConference(conferenceID string)
}
