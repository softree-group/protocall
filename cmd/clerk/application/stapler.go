package application

import (
	"context"
	"fmt"
	"net/url"
	"strings"
	"time"

	"protocall/api"
	"protocall/cmd/clerk/domain"
	"protocall/pkg/logger"
)

type Stapler struct {
	storage domain.RecordStorage
}

func NewStapler(s domain.RecordStorage) *Stapler {
	return &Stapler{
		storage: s,
	}
}

func (s *Stapler) collect(ctx context.Context, req *api.SendProtocolRequest) (string, error) {
	var res string
	for file := range s.storage.ListObjects(ctx, req.ConferenceID) {
		if strings.Contains(file.Key, ".wav") {
			raw, err := s.storage.GetFile(ctx, strings.Replace(file.Key, ".wav", ".txt", 1))
			if err != nil {
				return "", err
			}
			res += string(raw)
		}
	}
	
	return res, nil
}

func (s *Stapler) merge(ctx context.Context, req *api.SendProtocolRequest) (*url.URL, error) {
	var (
		protocol string
		err      error
	)

	const (
		retry = 3
		after = 10 * time.Second
	)

	for i := 0; i < retry; i++ {
		protocol, err = s.collect(ctx, req)
		if err != nil {
			logger.L.Error("error while collecting records", err)
			time.Sleep(after)
		}
	}
	if err != nil {
		return nil, err
	}

	path := fmt.Sprintf("%v/result.txt", req.ConferenceID)
	if err := s.storage.PutObject(
		ctx,
		path,
		strings.NewReader(protocol),
	); err != nil {
		return nil, err
	}

	res, err := s.storage.GetLink(ctx, path)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (s *Stapler) NewProtocol(
	req *api.SendProtocolRequest,
	send ...func(context.Context, string, []string),
) {
	go func() {
		ctx := context.Background()
		link, err := s.merge(ctx, req)
		if err != nil {
			logger.L.Error("error merge files for conference: ", req.ConferenceID)
			return
		}

		for _, job := range send {
			job(ctx, link.String(), req.To)
		}
		logger.L.Info("protocol was successfully send for: ", req.ConferenceID)
	}()
}
