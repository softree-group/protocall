package repository

import (
	"protocall/domain/entity"
)

type Record interface {
	RecordSender
	RecordStorage
	RecordTranslator
}

type RecordStorage interface {
	Upload(string, string) error
}

type RecordTranslator interface {
	Translate(*entity.User, *entity.Conference, string) error
}

type RecordSender interface {
	Send()
}
