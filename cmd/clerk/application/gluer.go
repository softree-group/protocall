package application

import (
	"context"
	"strings"

	"protocall/cmd/clerk/domain"
)

type Gluer struct {
	storage domain.RecordStorage
}

func NewGluer(s domain.RecordStorage) *Gluer {
	return &Gluer{
		storage: s,
	}
}

func (g *Gluer) Merge(ctx context.Context, conferenceID string) (string, error) {
	var res string
	for file := range g.storage.ListObjects(ctx, conferenceID) {
		if strings.Contains(file.Key, ".txt") {
			raw, err := g.storage.GetFile(ctx, file.Key)
			if err != nil {
				return "", err
			}
			res += string(raw)
		}
	}

	return res, nil
}
