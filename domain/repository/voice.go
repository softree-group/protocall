package repository

import (
	"context"
	"protocall/domain/entity"
)

type Voice interface {
	Recognize(context.Context, string) (*entity.Message, error)
}

type VoiceStorage interface {
	GetRecord(context.Context, string) ([]byte, error)
}

type VoiceRecognizer interface {
	Recognize(context.Context, []byte) (*entity.Message, error)
}
