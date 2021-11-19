package mailer

import (
	"context"
	"protocall/pkg/logger"

	"gopkg.in/gomail.v2"
)

type Mailer struct {
	client   *gomail.Dialer
	username string
}

func NewMailer(c *MailerConfig) *Mailer {
	return &Mailer{
		client:   gomail.NewDialer(c.Host, c.Port, c.Username, c.Password),
		username: c.Username,
	}
}

type Message struct {
	Type    string
	Subject string
	Body    string
}

func (m *Mailer) Send(ctx context.Context, msg *Message, to ...string) {
	for _, user := range to {
		newMsg := gomail.NewMessage()
		newMsg.SetHeader("From", m.username)
		newMsg.SetHeader("To", user)
		newMsg.SetHeader("Subject", msg.Subject)
		newMsg.SetBody(msg.Type, msg.Body)

		if err := m.client.DialAndSend(); err != nil {
			logger.L.Error("error while render template for user: ", user)
			continue
		}
	}
}
