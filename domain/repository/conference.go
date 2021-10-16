package repository

import "protocall/domain/entity"

type Conference interface {
	Get(conferenceID string) *entity.Conference
	Save(conference *entity.Conference)
	Delete(conferenceID string)
}
