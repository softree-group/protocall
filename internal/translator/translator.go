package translator

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"time"

	"protocall/internal/stapler"
	"protocall/pkg/logger"
)

var (
	errGetObject = errors.New("error while getting object from storage")
)

type Recognizer interface {
	Recognize(ctx context.Context, filename string, length time.Duration) (<-chan TranslateRespone, <-chan error)
}

type Storage interface {
	PutObject(context.Context, string, io.Reader) error
	GetObject(context.Context, string) (io.ReadCloser, error)
}

type Connector interface {
	TranslationDone(context.Context, *ConnectorRequest) error
}

type Translator struct {
	storage    Storage
	recognizer Recognizer
	connector  Connector
}

func NewTranslator(r Recognizer, s Storage, c Connector) *Translator {
	return &Translator{
		storage:    s,
		recognizer: r,
		connector:  c,
	}
}

func phrase(req *TranslateRequest, resp *TranslateRespone) string {
	if len(resp.Alternatives) == 0 {
		return ""
	}

	chunk := resp.Alternatives[0]
	if chunk.Text == "" {
		return ""
	}

	offset, err := time.ParseDuration(chunk.Words[0].StartTime)
	if err != nil {
		return ""
	}

	return fmt.Sprintf(
		"%v%s%v%s%v\n",
		req.ConnectTime.Add(offset).Format(time.RFC850),
		stapler.Delimeter,
		req.User.Username,
		stapler.Delimeter,
		chunk.Text,
	)
}

func (t *Translator) Translate(req *TranslateRequest) {
	go func() {
		ctx := context.Background()
		w := bytes.NewBuffer([]byte{})
		chunk, err := t.recognizer.Recognize(ctx, req.Record.URI, req.Record.Length)
		if err := func() error {
			for {
				select {
				case data := <-chunk:
					fmt.Fprint(w, phrase(req, &data))
				case err := <-err:
					return err
				}
			}
		}(); err != nil {
			logger.L.Error(err)
			return
		}

		if err := t.storage.PutObject(
			ctx,
			req.Text,
			bytes.NewReader(w.Bytes()),
		); err != nil {
			logger.L.Errorf("%w: %v", errGetObject, err)
			return
		}

		if err := t.connector.TranslationDone(context.Background(), &ConnectorRequest{
			SessionID: req.SessionID,
			Text:      req.Text,
			Record: Record{
				URI:  req.Record.URI,
				Path: req.Record.Path,
			},
		}); err != nil {
			logger.L.Errorln("error while notify connector: ", err)
			return
		}
		logger.L.Info("Connector successfully notified")
	}()
}
