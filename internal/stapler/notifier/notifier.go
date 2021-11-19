package notifier

import (
	"context"
	"protocall/internal/stapler"
)

type Runner interface {
	Send(context.Context, string, string, string, ...string)
}

type Notifier struct {
	runner []Runner
}

func NewNotifier(runners ...Runner) *Notifier {
	return &Notifier{runner: runners}
}

func (n *Notifier) Notify(ctx context.Context, phrases []stapler.Phrase, to ...string) {
	for _, el := range n.runner {
		el.Send(ctx, "text/html", subject, render(phrases), to...)
	}
}
