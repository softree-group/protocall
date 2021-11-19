package s3

import "os"

type StorageConfig struct {
	UseSSL    bool   `yaml:"useSSL"`
	Bucket    string `yaml:"bucket"`
	Endpoint  string `yaml:"endpoint"`
	AccessKey string
	SecretKey string
}

func ApplySecrets(cfg *StorageConfig) {
	cfg.AccessKey = os.Getenv("ACCESS_KEY")
	cfg.SecretKey = os.Getenv("SECRET_KEY")
}
