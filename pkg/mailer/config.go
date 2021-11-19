package mailer

import "os"

type MailerConfig struct {
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	Username string `yaml:"username"`
	Password string
}

func ApplySecrets(cfg *MailerConfig) {
	cfg.Password = os.Getenv("EMAIL_KEY")
}
