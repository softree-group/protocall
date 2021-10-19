package storage

import (
	"context"
	"io"
	"protocall/domain/repository"

	"github.com/minio/minio-go"
)

type Config struct {
	DisableSSL bool
	AccessKey  string
	SecretKey  string
	Endpoint   string
}

type s3client struct {
	mc *minio.Client
}

func NewStorage(c *Config) (repository.VoiceStorage, error) {
	mc, err := minio.New(c.Endpoint, c.AccessKey, c.SecretKey, !c.DisableSSL)
	if err != nil {
		return nil, err
	}

	return &s3client{mc: mc}, nil
}

func (c *s3client) UploadFile(ctx context.Context, bucketName, localPath, remotePath string) error {
	if _, err := c.mc.FPutObjectWithContext(
		ctx,
		bucketName,
		remotePath,
		localPath,
		minio.PutObjectOptions{},
	); err != nil {
		return err
	}

	return nil
}

func (c *s3client) GetFile(ctx context.Context, bucketName, remotePath string) (io.ReadCloser, error) {
	return c.mc.GetObjectWithContext(ctx, bucketName, remotePath, minio.GetObjectOptions{})
}
