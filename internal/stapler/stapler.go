package stapler

import (
	"context"
	"sort"

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

type Stapler struct {
	storage  Storage
	notifier Notifier
}

func NewStapler(storage Storage, notifier Notifier) *Stapler {
	return &Stapler{
		storage:  storage,
		notifier: notifier,
	}
}

func (s *Stapler) Protocol(ctx context.Context, req *ProtocolRequest) error {
	var data []Phrase
	for _, user := range req.Users {
		for _, text := range user.Texts {
			raw, err := s.storage.GetFile(ctx, text)
			if err != nil {
				return err
			}
			data = append(data, parseString(string(raw))...)
		}
	}

	sort.Slice(data, func(i, j int) bool {
		return data[i].Time.Unix() < data[j].Time.Unix()
	})

	s.notifier.Send(ctx, data, req.Users)
	logger.L.Info("successfully send protocol")
	return nil
}
