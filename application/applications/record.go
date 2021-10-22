package applications

import (
	"context"
	"protocall/domain/entity"
)

type Record interface {
	Translate(context.Context, string) (*entity.Message, error)
	UploadToStorage(context.Context) error
	SendToUser(context.Context)
}
