package app

import (
	"context"
	"protocall/domain/entity"
	"protocall/domain/repository"
)

type Record struct {
	reps repository.Repositories
}

func NewRecord(r repository.Repositories) *Record {
	return &Record{
		reps: r,
	}
}

func (r *Record) Translate(context.Context, string) (*entity.Message, error) {
	return nil, nil
}
func (r *Record) UploadToStorage(context.Context) error {
	return nil
}
func (r *Record) SendToUser(context.Context) {}
