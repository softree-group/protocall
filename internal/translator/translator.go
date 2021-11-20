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

type Translator struct {
	storage    Storage
	recognizer Recognizer
	job       *job
}

func NewTranslator(r Recognizer, s Storage) *Translator {
	return &Translator{
		storage:    s,
		recognizer: r,
		job:       newJob(),
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
			req.Start.Add(time.Duration(resp.Result[0].Start*float64(time.Second))),
			req.User.Username,
			resp.Text,
		)
	}

	w := bytes.NewBuffer([]byte{})
	for text := range t.recognizer.Recognize(ctx, audio) {
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
		t.job.create(req.User.Record)
		defer t.job.resolve(req.User.Record)

		if err := t.processAudio(context.Background(), req); err != nil {
			logger.L.Error("error while process record: ", req.User.Record)
			return
		}
		logger.L.Info("Translation done: ", req.User.Text)
	}()
}

// Track translation jobs and send signal, when all jobs are done.
func (t *Translator) Wait(records []string) <-chan struct{} {
	return t.job.wait(records)
}
