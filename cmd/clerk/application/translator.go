package application

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"strings"
	"sync"
	"time"

	"protocall/api"
	"protocall/cmd/clerk/domain"
	"protocall/pkg/logger"

	"github.com/hashicorp/go-uuid"
)

const (
	Running = iota
	Ready
	Failed
)

var (
	errGetObject = errors.New("error while getting object from storage")
)

type Translator struct {
	storage    domain.RecordStorage
	recognizer domain.Recognizer

	jobStatus *sync.Map
}

func NewTranslator(r domain.Recognizer, s domain.RecordStorage) *Translator {
	return &Translator{
		storage:    s,
		recognizer: r,
		jobStatus:  new(sync.Map),
	}
}

func (t *Translator) GetStatus(jobID string) (int, error) {
	status, ok := t.jobStatus.Load(jobID)
	if !ok {
		return -1, fmt.Errorf("%v", "job not found")
	}
	return status.(int), nil
}

func phrase(conferenceStart *time.Time, resp *api.TextRespone) string {
	if resp.Text == "" {
		return ""
	}
	return fmt.Sprintf(
		"%v:%v\n",
		conferenceStart.Add(time.Duration(resp.Result[0].Start*float64(time.Second))),
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
		fmt.Fprint(w, phrase(&req.StartTime, &text))
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

func (t *Translator) CreateJob(req *api.TranslateRequest) (string, error) {
	jobID, err := uuid.GenerateUUID()
	if err != nil {
		return "", err
	}
	t.jobStatus.Store(jobID, Running)

	go func() {
		if err := t.processAudio(context.Background(), req); err != nil {
			t.jobStatus.Store(jobID, Failed)
			return
		}
		t.jobStatus.Store(jobID, Ready)
		logger.L.Info("Translation done ", req.User.Path)
	}()

	return jobID, nil
}
