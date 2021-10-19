package recognizer

import (
	"context"
	"protocall/domain/entity"
	"protocall/domain/repository"
)

type RecognizerConfig struct {
	Host string
	Port string
}

type recognizer struct{}

func NewRecognizer(c *RecognizerConfig) repository.VoiceRecognizer {
	return &recognizer{}
}

func (r *recognizer) Recognize(context.Context, []byte) (*entity.Message, error) {
	return nil, nil
}
