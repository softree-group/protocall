package application

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"protocall/api"
	"protocall/cmd/clerk/domain"
	"protocall/pkg/logger"
)

var (
	errGetObject = errors.New("error while getting object from storage")
)

type Translator struct {
	storage    domain.RecordStorage
	recognizer domain.Recognizer
}

func NewTranslator(r domain.Recognizer, s domain.RecordStorage) *Translator {
	return &Translator{
		storage:    s,
		recognizer: r,
	}
}

func phrase(req *api.TranslateRequest, resp *api.TextRespone) string {
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

func (t *Translator) processAudio(ctx context.Context, req *api.TranslateRequest) error {
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

func (t *Translator) TranslateRecord(req *api.TranslateRequest) {
	go func() {
		if err := t.processAudio(context.Background(), req); err != nil {
			logger.L.Error("error while process record: ", req.User.Path)
			return
		}
		logger.L.Info("Translation done ", req.User.Path)
	}()
}
