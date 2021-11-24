package translator

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"time"

	"protocall/internal/stapler"
	"protocall/pkg/logger"
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

func (t *Translator) Translate(req *TranslateRequest) {
	go func() {
		ctx := context.Background()
		w := bytes.NewBuffer([]byte{})
		chunk, err := t.recognizer.Recognize(ctx, req.Record.URI, req.Record.Length)
		if err := func() error {
			for {
				select {
				case data := <-chunk:
					if len(data.Alternatives) == 0 {
						return nil
					}

					chunk := data.Alternatives[0]
					if chunk.Text == "" {
						return nil
					}

					offset, err := time.ParseDuration(chunk.Words[0].StartTime)
					if err != nil {
						return err
					}

					fmt.Fprintf(w, "%v%s%v%s%v\n",
						req.ConnectTime.Add(offset).Format(time.RFC850),
						stapler.Delimeter,
						req.User.Username,
						stapler.Delimeter,
						chunk.Text,
					)
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
			logger.L.Errorf("%w: %v", "error while getting object from storage", err)
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
