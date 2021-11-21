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
	Recognize(context.Context, io.Reader) <-chan TranslateRespone
}

type Storage interface {
	PutObject(context.Context, string, io.Reader) error
	GetObject(context.Context, string) (io.ReadCloser, error)
}

type Connector interface {
	TranslationDone(context.Context, *TranslateRequest) error
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
	audio, err := t.storage.GetObject(ctx, req.User.Record)
	if err != nil {
		return err
	}
	defer audio.Close()

	phrase := func(req *TranslateRequest, resp *TranslateRespone) string {
		if resp.Text == "" {
			return ""
		}

		return fmt.Sprintf(
			"%v:%v:%v\n",
			req.ConnectTime.Add(time.Duration(resp.Result[0].Start*float64(time.Second))),
			req.User.Username,
			resp.Text,
		)
	}

	w := bytes.NewBuffer([]byte{})
	for text := range t.recognizer.Recognize(ctx, audio) {
		fmt.Println(text)
		fmt.Fprint(w, phrase(req, &text))
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
			logger.L.Errorln("error while process record: ", req.User.Record)
			return
		}
		logger.L.Info("Translation done: ", req.User.Text)
		if err := t.connector.TranslationDone(context.Background(), req); err != nil {
			logger.L.Errorln("error while notify connector: ", err)
			return
		}
		logger.L.Info("Connector successfully notified: ", req.User.SessionID)
	}()
}
