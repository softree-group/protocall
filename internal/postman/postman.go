package postman

import (
	"context"

	"protocall/internal/stapler"
	"protocall/pkg/logger"
)

type Runner interface {
	Send(context.Context, string, string, string, string) error
}

type Postman struct {
	mail Runner
}

func NewPostman(mail Runner) *Postman {
	return &Postman{mail}
}

func (p *Postman) Send(ctx context.Context, protocol []stapler.Phrase, users []stapler.User) {
	payload := render(protocol)
	for _, user := range users {
		if user.NeedProtocol {
			err := p.mail.Send(ctx, "text/html", subject, payload, user.Email)
			if err != nil {
				logger.L.Errorln(err)
			}
		}
	}
	logger.L.Info("Protocol successfully sended")
}
