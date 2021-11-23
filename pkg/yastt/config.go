package yastt

import "os"

type YasttConfig struct {
	TranscribeAddr  string `yaml:"transcribeAddr"`
	OperationAddr   string `yaml:"operationAddr"`
	Specification   `yaml:"specification"`
	PoolCoefficient float64 `yaml:"poolCoefficient"`
	Token           string
}

func ApplySecrets(cfg *YasttConfig) {
	cfg.Token = os.Getenv("RECOGNIZER_KEY")
}
