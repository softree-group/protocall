package translator

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"strings"
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

type Translator struct {
	storage    Storage
	recognizer Recognizer
}

func NewTranslator(r Recognizer, s Storage) *Translator {
	return &Translator{
		storage:    s,
		recognizer: r,
	}
}

func phrase(req *TranslateRequest, resp *TranslateRespone) string {
	if resp.Text == "" {
		return ""
	}
	return fmt.Sprintf(
		"%v:%v:%v\n",
		req.StartTime.Add(time.Duration(resp.Result[0].Start*float64(time.Second))),
		req.User.Username,
		resp.Text,
	)
}

func (t *Translator) processAudio(ctx context.Context, req *TranslateRequest) error {
	audio, err := t.storage.GetObject(ctx, req.User.Path)
	if err != nil {
		return err
	}
	defer audio.Close()

	w := bytes.NewBuffer([]byte{})
	for text := range t.recognizer.Recognize(ctx, audio) {
		fmt.Fprint(w, phrase(req, &text))
	}

	if err := t.storage.PutObject(
		ctx,
		strings.Replace(req.User.Path, ".wav", ".txt", 1),
		bytes.NewReader(w.Bytes()),
	); err != nil {
		return fmt.Errorf("%w: %v", errGetObject, err)
	}

	return nil
}

func (t *Translator) Translate(req *TranslateRequest) {
	go func() {
		if err := t.processAudio(context.Background(), req); err != nil {
			logger.L.Error("error while process record: ", req.User.Path)
			return
		}
		logger.L.Info("Translation done ", req.User.Path)
	}()
}
