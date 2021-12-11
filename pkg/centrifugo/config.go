package centrifugo

import "os"


CentrifugoToken  = "centrifugo.token"
CentrifugoAPIKey = "centrifugo.api_key"
CentrifugoHost   = "centrifugo.host"

type Config struct {
	Addr  string `yaml:"addr"`
	Token string
}

func applySecrets(cfg *Config) {
	cfg.Token = os.Getenv("CENTRIFUGO_KEY")
}
