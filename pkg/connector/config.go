package connector

import "os"

type ConnectorClientConfig struct {
	Host  string `yaml:"host"`
	Port  string `yaml:"port"`
	Token string
}

func ApplySecrets(c *ConnectorClientConfig) {
	c.Token = os.Getenv("CONNECTOR_KEY")
}
