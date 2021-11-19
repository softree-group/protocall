package stapler

import (
	"bufio"
	"context"
	"fmt"
	"net/url"
	"sort"
	"strings"
	"time"

	"protocall/api"
	"protocall/internal/clerk/repositories"
	"protocall/pkg/logger"
)

type Stapler struct {
	storage repositories.Storage
}

func NewStapler(s repositories.Storage) *Stapler {
	return &Stapler{
		storage: s,
	}
}

type phraseLine struct {
	time time.Time
	user string
	text string
}

func (s *Stapler) collect(ctx context.Context, req *api.SendProtocolRequest) ([]phraseLine, error) {
	var rawRes string
	for file := range s.storage.ListObjects(ctx, req.ConferenceID) {
		if strings.Contains(file.Key, ".wav") {
			raw, err := s.storage.GetFile(ctx, strings.Replace(file.Key, ".wav", ".txt", 1))
			if err != nil {
				return nil, err
			}
			rawRes += string(raw)
		}
	}

	scanner := bufio.NewScanner(strings.NewReader(rawRes))
	scanner.Split(bufio.ScanLines)
	res := []phraseLine{}
	for scanner.Scan() {
		line := strings.Split(scanner.Text(), ":")
		if len(line) < 2 {
			return nil, fmt.Errorf("invalid line")
		}

		time, err := time.Parse(time.RFC3339, line[0])
		if err != nil {
			return nil, err
		}

		res = append(res, phraseLine{
			time: time,
			user: line[1],
			text: line[2],
		})
	}

	sort.Slice(res, func(i, j int) bool {
		return res[i].time.Unix() < res[j].time.Unix()
	})

	return res, nil
}

func (s *Stapler) merge(ctx context.Context, req *api.SendProtocolRequest) (*url.URL, error) {
	const (
		retry = 3
		after = 10 * time.Second
	)

	var (
		err error
	)

	for i := 0; i < retry; i++ {
		lines, err = s.collect(ctx, req)
		if err != nil {
			logger.L.Error("error while collecting records", err)
			time.Sleep(after)
		}
	}
	if err != nil {
		return nil, err
	}

	// path := fmt.Sprintf("%v/result.txt", req.ConferenceID)
	// if err := s.storage.PutObject(
	// 	ctx,
	// 	path,
	// 	strings.NewReader(protocol),
	// ); err != nil {
	// 	return nil, err
	// }

	// res, err := s.storage.GetLink(ctx, path)
	// if err != nil {
	// 	return nil, err
	// }

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
