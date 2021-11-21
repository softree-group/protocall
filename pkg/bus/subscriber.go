package bus

import "github.com/sirupsen/logrus"

type Subscriber struct {
	C     chan interface{}
	clear func()
	event string
	uid string
}

func (s *Subscriber) Cancel() {
	logrus.Warn("CANCEL ", s.event, " ", s.uid)
	close(s.C)
	s.C = nil

	s.clear()
}

func (s Subscriber) Channel() chan interface{} {
	return s.C
}

var _ *Subscriber = &Subscriber{}
