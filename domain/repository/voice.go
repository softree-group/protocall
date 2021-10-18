package repository

import (
	"context"
	"protocall/domain/entity"
	"protocall/infrastructure/storage"
)

type Voice interface {
	Recognize(context.Context, string) (*entity.Message, error)
}

type VoiceStorage interface {
	UploadFile(bucketName string, remotePath string, localPath string) (string, error)
	DownloadFile(bucketName string, remotePath string, versionID string, localPath string) error
}

type VoiceRecognizer interface {
	Recognize(context.Context, []byte) (*entity.Message, error)
}
