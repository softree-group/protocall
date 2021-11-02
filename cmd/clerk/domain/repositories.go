package domain

import (
	"context"
	"io"

	"protocall/api"
	"protocall/pkg/s3"
)

type Gluer interface {
	Merge(context.Context, string) (string, error)
}

type Sender interface {
	SendProtocol(context.Context, string, []string) error
}

type Translator interface {
	GetStatus(string) (int, error)
	CreateJob(*api.TranslateRequest) (string, error)
}

type RecordStorage interface {
	PutFile(context.Context, string, string) error
	PutObject(context.Context, string, io.Reader) error
	GetObject(context.Context, string) (io.ReadCloser, error)
	GetFile(context.Context, string) ([]byte, error)
	ListObjects(context.Context, string) <-chan s3.ObjectInfo
}

type Recognizer interface {
	Recognize(context.Context, io.Reader) <-chan api.TextRespone
}
