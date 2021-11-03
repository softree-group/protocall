package logger

import (
	"os"

	"github.com/sirupsen/logrus"
)

type LoggerConfig struct {
	OutputPath string `yaml:"logOutput"`
	LogLevel   string `yaml:"logLevel"`
}

var L *logrus.Logger

func NewLogger(c *LoggerConfig) error {
	L = logrus.New()

	L.SetOutput(os.Stdout)
	if c.OutputPath != "" {
		f, err := os.Open(c.OutputPath)
		if err != nil {
			return err
		}
		L.SetOutput(f)
	}

	level, err := logrus.ParseLevel(c.LogLevel)
	if err != nil {
		return err
	}
	L.SetLevel(level)
	return nil
}
