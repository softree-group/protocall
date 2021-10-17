package storage

import (
	"context"
	"io"
	"os"
	"protocall/domain/repository"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"

	"protocall/domain/repository"
)

type Config struct {
	Endpoint string
	Bucket   string
	KeyId    string
	Key      string
}

type s3client struct {
	client  *s3.S3
	session *session.Session

	progressOutput io.Writer
}

type UploadFileOptions struct {
	Acl                  string
	ServerSideEncryption string
	KmsKeyId             string
	ContentType          string
	DisableMultipart     bool
}

func NewStorage(c *Config) (repository.VoiceStorage, error) {
	sess := session.Must(session.NewSession(&aws.Config{
		Region:   aws.String("us-west-2"),
		Endpoint: &c.Endpoint,
	}))

	// creds := stscreds.NewCredentials(sess, "recognizer")

	return &Storage{
		sess:   sess,
		config: c,
	}, nil
}

func (s *Storage) Download(ctx context.Context, filename string) ([]byte, error) {
	d := s3manager.NewDownloader(s.sess)
	res := aws.NewWriteAtBuffer([]byte{})
	_, err := d.DownloadWithContext(
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
	return res.Bytes(), nil
}

func (s *Storage) Upload(ctx context.Context, file *os.File) error {
	u := s3manager.NewUploader(s.sess)
	if _, err := u.UploadWithContext(ctx, &s3manager.UploadInput{
		Bucket: aws.String(s.config.Bucket),
		Key:    aws.String(file.Name()),
		Body:   file,
	}); err != nil {
		return err
	}
	return nil
	return res, nil
}
