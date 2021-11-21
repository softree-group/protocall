package repository

import (
	"context"
	"protocall/internal/connector/domain/entity"
	"protocall/internal/stapler"
	"protocall/internal/translator"
)

type Conference interface {
	GetConference(conferenceID string) *entity.Conference
	SaveConference(conference *entity.Conference)
	DeleteConference(conferenceID string)
}

type ConferenceStorage interface {
	UploadRecord(ctx context.Context, path string) error
}

type ConferenceTranslator interface {
	TranslateRecord(ctx context.Context, data *translator.TranslateRequest) error
	CreateProtocol(ctx context.Context, data *stapler.ProtocolRequest) error
}

type ConferenceJobs interface {
	Store(conferenceID, recName string) error
	DoneJob(conferenceID, recName string) error
	IsDone(conferenceID string) (bool, error)
}
