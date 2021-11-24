package stapler

import (
	"context"
	"sort"

	"protocall/pkg/s3"
)

type Storage interface {
	ListObjects(context.Context, string) <-chan s3.ObjectInfo
	GetFile(context.Context, string) ([]byte, error)
}

type Stapler struct {
	storage Storage
}

func NewStapler(storage Storage) *Stapler {
	return &Stapler{storage}
}

func (s *Stapler) Make(ctx context.Context, req *ProtocolRequest) ([]Phrase, error) {
	var data []Phrase
	for _, user := range req.Users {
		for _, text := range user.Texts {
			raw, err := s.storage.GetFile(ctx, text)
			if err != nil {
				return nil, err
			}

			parsed, err := newPhrase(string(raw))
			if err != nil {
				return nil, err
			}
			data = append(data, parsed...)
		}
	}

	sort.Slice(data, func(i, j int) bool {
		return data[i].Time.Unix() < data[j].Time.Unix()
	})
	return data, nil
}
