package storage

import (
	"context"
	"io"
	"protocall/domain/repository"

	"github.com/minio/minio-go"
)

type StorageConfig struct {
	DisableSSL bool   `yaml:"disableSSL"`
	Bucket     string `yaml:"bucket"`
	Endpoint   string `yaml:"endpoint"`
	AccessKey  string
	SecretKey  string
}

type s3client struct {
	client *minio.Client
}

func NewStorage(c *StorageConfig) (repository.VoiceStorage, error) {
	mc, err := minio.New(c.Endpoint, c.AccessKey, c.SecretKey, !c.DisableSSL)
	if err != nil {
		return nil, err
	}

	return &s3client{client: mc}, nil
}

func (c *s3client) UploadFile(ctx context.Context, bucketName, localPath, remotePath string) error {
	if _, err := c.client.FPutObjectWithContext(
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
	return c.client.GetObjectWithContext(ctx, bucketName, remotePath, minio.GetObjectOptions{})
}
