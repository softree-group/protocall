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
	runners []Runner
}

func NewNotifier(runners ...Runner) *Notifier {
	return &Notifier{runners: runners}
}

func (n *Notifier) Notify(ctx context.Context, phrases []stapler.Phrase, to ...string) {
	for _, runner := range n.runners {
		for _, user := range to {
			err := runner.Send(ctx, "text/html", subject, render(phrases), user)
			if err != nil {
				logger.L.Errorln(err)
			}
		}
	}
}
