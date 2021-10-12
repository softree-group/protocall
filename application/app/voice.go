package app

import (
	"context"
	"protocall/domain/entity"
	"protocall/domain/repository"
)

type Voice struct {
	storage repository.VoiceStorage
	tts     repository.VoiceRecognizer
}

func NewVoice(s repository.VoiceStorage, r repository.VoiceRecognizer) repository.Voice {
	return &Voice{
		s, r,
	}
}

func (v *Voice) Recognize(ctx context.Context, filename string) (*entity.Message, error) {
	audio, err := v.storage.GetRecord(ctx, filename)
	if err != nil {
		return nil, err
	}
	msg, err := v.tts.Recognize(ctx, audio)
	if err != nil {
		return nil, err
	}
	return msg, nil
}
