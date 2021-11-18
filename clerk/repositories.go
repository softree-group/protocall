package clerk

import (
	"context"
	"io"
	"net/url"

	"protocall/api"
	"protocall/pkg/s3"
)

type Stapler interface {
	NewProtocol(
		*api.SendProtocolRequest,
		...func(context.Context, string, []string),
	)
}

type Sender interface {
	SendSMTP(context.Context, string, []string)
}

type Translator interface {
	TranslateRecord(*api.TranslateRequest)
}

type RecordStorage interface {
	GetLink(ctx context.Context, path string) (*url.URL, error)
	PutFile(context.Context, string, string) error
	PutObject(context.Context, string, io.Reader) error
	GetObject(context.Context, string) (io.ReadCloser, error)
	GetFile(context.Context, string) ([]byte, error)
	ListObjects(context.Context, string) <-chan s3.ObjectInfo
}

type Recognizer interface {
	Recognize(context.Context, io.Reader) <-chan api.TextRespone
}
