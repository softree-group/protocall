package notifier

import (
	"context"

	"protocall/internal/stapler"
	"protocall/pkg/logger"
)

type Runner interface {
	Send(context.Context, string, string, string, string) error
}

type Notifier struct {
	mail Runner
}

func NewNotifier(mail Runner) *Notifier {
	return &Notifier{mail: mail}
}

func (n *Notifier) Send(ctx context.Context, protocol []stapler.Phrase, users []stapler.User) {
	payload := render(protocol)
	for _, user := range users {
		if user.NeedProtocol {
			err := n.mail.Send(ctx, "text/html", subject, payload, user.Email)
			if err != nil {
				logger.L.Errorln(err)
			}
		}
	}
	logger.L.Info("Protocol successfully sended")
}
