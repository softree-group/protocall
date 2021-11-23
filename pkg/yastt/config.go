package yastt

import "os"

type YasttConfig struct {
	TranscribeAddr string `yaml:"transcribeAddr"`
	OperationAddr  string `yaml:"operationAddr"`
	Specification         `yaml:"specification"`
	// 1 minute of single channel audio is recognized in about 10 seconds.
	// PoolCoefficient: 1.7
	PoolCoefficient float64 `yaml:"poolCoefficient"`
	AccessKey       string
	SecretKey       string
}

func ApplySecrets(cfg *YasttConfig) {
	cfg.AccessKey = os.Getenv("YASTT_ACCESS_KEY")
	cfg.SecretKey = os.Getenv("YASTT_SECRET_KEY")
}
