package application

import (
	"bytes"
	"context"
	"html/template"

	"protocall/pkg/logger"

	"gopkg.in/gomail.v2"
)

var (
	smtpTmpl *template.Template
)

func init() {
	smtpTmpl = template.Must(template.New("base").Parse(protocallTmpl))
}

type SenderConfig struct {
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	Username string `yaml:"username"`
	Password string
}

type Sender struct {
	client *gomail.Dialer
	config *SenderConfig
}

func NewSender(c *SenderConfig) *Sender {
	return &Sender{
		client: gomail.NewDialer(c.Host, c.Port, c.Username, c.Password),
		config: c,
	}
}

const protocallTmpl = `<p>По этой <a href="{{.Link}}">ссылке</a> вы можете скачать стенограмму конференции.</p>`

func (s *Sender) SendSMTP(ctx context.Context, linkToDownload string, users []string) {
	for _, user := range users {
		msg := bytes.Buffer{}
		if err := smtpTmpl.Execute(
			&msg,
			struct {
				Link template.URL
			}{
				Link: template.URL(linkToDownload),
			},
		); err != nil {
			logger.L.Error("error while render template for user: ", user)
			continue
		}

		m := gomail.NewMessage()
		m.SetHeader("From", s.config.Username)
		m.SetHeader("To", user)
		m.SetHeader("Subject", "Запись вашей конференции")
		m.SetBody("text/html", msg.String())

		if err := s.client.DialAndSend(m); err != nil {
			logger.L.Error("error while render template for user: ", user)
			continue
		}
	}
}
