package storage

import (
	"context"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"

	"protocall/domain/repository"
)

type Config struct {
	Endpoint string
	Bucket   string
}

type Storage struct {
	downloader *s3manager.Downloader
	config     *Config
}

func NewStorage(c *Config) repository.VoiceStorage {
	return &Storage{
		downloader: s3manager.NewDownloader(
			s3.New(session.New(&aws.Config{})),
		),
		config: c,
	}
}

func (s *Storage) GetRecord(ctx context.Context, filename string) ([]byte, error) {
	var res []byte
	_, err := s.downloader.DownloadWithContext(
		ctx,
		res,
		&s3.GetObjectInput{
			Bucket: aws.String(s.config.Bucket),
			Key:    aws.String(filename),
		},
	)
	if err != nil {
		return nil, err
	}
	return res, nil
}
