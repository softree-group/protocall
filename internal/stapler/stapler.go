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
	Notify(context.Context, []Phrase, ...string)
}

type Stapler struct {
	storage  Storage
	notifier Notifier
}

func NewStapler(s Storage, n Notifier) *Stapler {
	return &Stapler{
		storage:  s,
		notifier: n,
	}
}

func (s *Stapler) collect(ctx context.Context, req *ProtocolRequest) ([]Phrase, error) {
	var lines string
	for file := range s.storage.ListObjects(ctx, req.ConferenceID) {
		if strings.Contains(file.Key, ".wav") {
			raw, err := s.storage.GetFile(ctx, strings.Replace(file.Key, ".wav", ".txt", 1))
			if err != nil {
				return nil, err
			}
			lines += string(raw)
		}
	}

	scanner := bufio.NewScanner(strings.NewReader(lines))
	scanner.Split(bufio.ScanLines)
	res := []Phrase{}
	for scanner.Scan() {
		line := strings.Split(scanner.Text(), ":")
		if len(line) < 2 {
			return nil, fmt.Errorf("invalid line")
		}

		time, err := time.Parse(time.RFC3339, line[0])
		if err != nil {
			return nil, err
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

const (
	retry = 3
	after = 10 * time.Second
)

func (s *Stapler) Protocol(req *ProtocolRequest) {
	go func() {
		var (
			phrases []Phrase
			err     error
		)

		ctx := context.Background()

		for i := 0; i < retry; i++ {
			phrases, err = s.collect(ctx, req)
			if err != nil {
				time.Sleep(after)
			}
		}
		if err != nil {
			logger.L.Error("error while collecting records", err)
			return
		}

		s.notifier.Notify(ctx, phrases, req.To...)
		logger.L.Info("successfully send protocol for conference: ", req.ConferenceID)
	}()
}
