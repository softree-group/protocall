package s3

import (
	"context"
	"fmt"
	"io"
	"net/url"
	"strings"
	"time"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

const (
	unkownSize = -1
	expire     = time.Second * 24 * 60 * 60
)

type StorageConfig struct {
	UseSSL    bool   `yaml:"useSSL"`
	Bucket    string `yaml:"bucket"`
	Endpoint  string `yaml:"endpoint"`
	AccessKey string
	SecretKey string
}

type S3 struct {
	client *minio.Client
	bucket string
}

func NewStorage(c *StorageConfig) (*S3, error) {
	mc, err := minio.New(c.Endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(c.AccessKey, c.SecretKey, ""),
		Secure: c.UseSSL,
	})
	if err != nil {
		return nil, err
	}

	return &S3{
		client: mc,
		bucket: c.Bucket,
	}, nil
}

func (s *S3) PutFile(ctx context.Context, localPath, remotePath string) error {
	if _, err := s.client.FPutObject(
		ctx,
		s.bucket,
		remotePath,
		localPath,
		minio.PutObjectOptions{},
	); err != nil {
		return err
	}

	return nil
}

func (s *S3) PutObject(ctx context.Context, objectName string, input io.Reader) error {
	if _, err := s.client.PutObject(
		ctx,
		s.bucket,
		objectName,
		input,
		unkownSize,
		minio.PutObjectOptions{},
	); err != nil {
		return err
	}
	return nil
}

func (s *S3) GetObject(ctx context.Context, path string) (io.ReadCloser, error) {
	return s.client.GetObject(ctx, s.bucket, path, minio.GetObjectOptions{})
}

func (s *S3) GetFile(ctx context.Context, path string) ([]byte, error) {
	file, err := s.client.GetObject(ctx, s.bucket, path, minio.GetObjectOptions{})
	if err != nil {
		return nil, err
	}
	defer file.Close()

	data, err := io.ReadAll(file)
	if err != nil {
		return nil, err
	}
	return data, nil
}

type ObjectInfo = minio.ObjectInfo

func (s *S3) ListObjects(ctx context.Context, path string) <-chan ObjectInfo {
	return s.client.ListObjects(
		ctx,
		s.bucket,
		minio.ListObjectsOptions{
			Prefix:    path,
			Recursive: true,
		},
	)
}

func (s *S3) GetLink(ctx context.Context, path string) (*url.URL, error) {
	reqParams := make(url.Values)
	reqParams.Set(
		"response-content-disposition",
		fmt.Sprintf("attachment; filename=%v", path[strings.LastIndex(path, "/")+1:]),
	)

	return s.client.PresignedGetObject(
		ctx,
		s.bucket,
		path,
		expire,
		reqParams,
	)
}
