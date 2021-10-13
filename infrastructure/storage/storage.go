package storage

import (
	"context"
	"io"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials/stscreds"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
)

type Config struct {
	Endpoint string
	Bucket   string
	KeyId    string
	Key      string
}

type Storage struct {
	sess   *session.Session
	config *Config
}

func NewStorage(c *Config) (*Storage, error) {
	sess := session.Must(session.NewSession(&aws.Config{
		Region:   aws.String("us-west-2"),
		Endpoint: &c.Endpoint,
	}))

	creds := stscreds.NewCredentials(sess, "recognizer")

	return &Storage{
		sess:   sess,
		config: c,
	}, nil
}

func (s *Storage) GetRecord(ctx context.Context, filename string) ([]byte, error) {
	d := s3manager.NewDownloader(s.sess)
	res := aws.NewWriteAtBuffer([]byte{})
	_, err := d.DownloadWithContext(
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
	return res.Bytes(), nil
}

func (s *Storage) UploadRecord(ctx context.Context, filename string, file io.Reader) error {
	u := s3manager.NewUploader(s.sess)
	if _, err := u.UploadWithContext(ctx, &s3manager.UploadInput{
		Bucket: aws.String(s.config.Bucket),
		Key:    aws.String(filename),
		Body:   file,
	}); err != nil {
		return err
	}
	return nil
}
