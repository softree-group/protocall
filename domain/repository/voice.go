package repository

import (
	"context"
	"io"

	"protocall/domain/entity"
)

type Voice interface {
	Translate(context.Context, string) (*entity.Message, error)
	SendToUser(context.Context)
}

type VoiceStorage interface {
	UploadFile(context.Context, string, string, string) error
	GetFile(context.Context, string, string) (io.ReadCloser, error)
}

type VoiceRecognizer interface {
	Recognize(context.Context, []byte) (*entity.Message, error)
}

type VoiceSender interface {
	SendSMTP()
}
