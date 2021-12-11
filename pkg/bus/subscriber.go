package bus

type Subscriber struct {
	C     chan interface{}
	clear func()
	event string
	uid   string
}

func (s *Subscriber) Cancel() {
	close(s.C)
	s.C = nil

	s.clear()
}

func (s Subscriber) Channel() chan interface{} {
	return s.C
}
