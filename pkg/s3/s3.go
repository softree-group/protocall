package s3

import (
	"context"
	"io"

	"github.com/minio/minio-go"
)

type StorageConfig struct {
	DisableSSL bool   `yaml:"disableSSL"`
	Bucket     string `yaml:"bucket"`
	Endpoint   string `yaml:"endpoint"`
	AccessKey  string
	SecretKey  string
}

type S3 struct {
	client *minio.Client
}

func NewStorage(c *StorageConfig) (*S3, error) {
	mc, err := minio.New(c.Endpoint, c.AccessKey, c.SecretKey, !c.DisableSSL)
	if err != nil {
		return nil, err
	}

	return &S3{client: mc}, nil
}

func (c *S3) UploadFile(ctx context.Context, bucketName, localPath, remotePath string) error {
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

func (c *S3) GetFile(ctx context.Context, bucketName, remotePath string) (io.ReadCloser, error) {
	return c.client.GetObjectWithContext(ctx, bucketName, remotePath, minio.GetObjectOptions{})
}
