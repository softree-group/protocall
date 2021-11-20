package stapler

import (
	"bufio"
	"context"
	"fmt"
	"sort"
	"strings"
	"time"

	"protocall/pkg/logger"
	"protocall/pkg/s3"
)

type Storage interface {
	ListObjects(context.Context, string) <-chan s3.ObjectInfo
	GetFile(context.Context, string) ([]byte, error)
}

type Notifier interface {
	Send(context.Context, []Phrase, []User)
}

type Translator interface {
	Wait([]string) <-chan struct{}
}

type Stapler struct {
	storage    Storage
	notifier   Notifier
	translator Translator
}

func NewStapler(storage Storage, notifier Notifier, translator Translator) *Stapler {
	return &Stapler{
		storage:    storage,
		notifier:   notifier,
		translator: translator,
	}
}

func (s *Stapler) collect(ctx context.Context, users []User) ([]Phrase, error) {
	var data string
	for _, user := range users {
		raw, err := s.storage.GetFile(ctx, user.Text)
		if err != nil {
			return nil, err
		}
		data += string(raw)

	}

	scanner := bufio.NewScanner(strings.NewReader(data))
	scanner.Split(bufio.ScanLines)
	res := []Phrase{}
	for scanner.Scan() {
		line := strings.Split(scanner.Text(), ":")
		if len(line) < 2 {
			return nil, fmt.Errorf("invalid line")
		}

		time, err := time.Parse(time.RFC3339, line[0])
		if err != nil {
			return nil, fmt.Errorf("invalid timestamp: %w", err)
		}

		res = append(res, Phrase{
			Time: time,
			User: line[1],
			Text: line[2],
		})
	}

	sort.Slice(res, func(i, j int) bool {
		return res[i].Time.Unix() < res[j].Time.Unix()
	})

	return res, nil
}

func (s *Stapler) Protocol(ctx context.Context, req *ProtocolRequest) error {
	<-s.translator.Wait(func() []string {
		records := make([]string, len(req.Users))
		for _, user := range req.Users {
			records = append(records, user.Record)
		}
		return records
	}())

	protocol, err := s.collect(ctx, req.Users)
	if err != nil {
		return err
	}

	s.notifier.Send(ctx, protocol, req.Users)
	logger.L.Info("successfully send protocol")
	return nil
}
