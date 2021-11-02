package repository

import "protocall/domain/entity"

type Conference interface {
	GetConference(conferenceID string) *entity.Conference
	SaveConference(conference *entity.Conference)
	DeleteConference(conferenceID string)
}

type ConferenceStorage interface {
	UploadConference(path string) error
}

type ConferenceTranslator interface {
	TranslateConference(user *entity.User, conference *entity.Conference) error
}
