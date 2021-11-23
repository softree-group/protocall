package translator

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"time"

	"protocall/pkg/logger"
)

var (
	errGetObject = errors.New("error while getting object from storage")
)

type Recognizer interface {
	Recognize(ctx context.Context, filename string, length time.Duration) <-chan TranslateRespone
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

func (t *Translator) processAudio(ctx context.Context, req *TranslateRequest) error {
	phrase := func(req *TranslateRequest, resp *TranslateRespone) string {
		if len(resp.Alternatives) == 0 {
			return ""
		}

		chunk := resp.Alternatives[0]
		if chunk.Text == "" {
			return ""
		}

		return fmt.Sprintf(
			"%v:%v:%v\n",
			req.ConnectTime.Add(chunk.Words[0].StartTime),
			req.User.Username,
			chunk.Text,
		)
	}

	w := bytes.NewBuffer([]byte{})
	for chunk := range t.recognizer.Recognize(ctx, req.Record.Path, req.Record.Length) {
		fmt.Fprint(w, phrase(req, &chunk))
	}

	if err := t.storage.PutObject(
		ctx,
		req.User.Text,
		bytes.NewReader(w.Bytes()),
	); err != nil {
		return fmt.Errorf("%w: %v", errGetObject, err)
	}

	return nil
}

func (t *Translator) Translate(req *TranslateRequest) {
	go func() {
		if err := t.processAudio(context.Background(), req); err != nil {
			logger.L.Errorln("error while process request: ", req)
			return
		}

		logger.L.Info("Translation done: ", req.Text)
		if err := t.connector.TranslationDone(context.Background(), &ConnectorRequest{
			SessionID: req.SessionID,
			Text:      req.Text,
			Record:    req.Record.Path,
		}); err != nil {
			logger.L.Errorln("error while notify connector: ", err)
			return
		}
		logger.L.Info("Connector successfully notified: ")
	}()
}
