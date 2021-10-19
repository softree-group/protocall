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

func (v *Voice) Translate(context.Context, string) (*entity.Message, error) {
	return nil, nil
}

func (v *Voice) SendToUser(context.Context) {
	return
}
