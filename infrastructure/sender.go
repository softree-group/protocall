package infrastructure

type Sender struct{}

func NewSender() *Sender {
	return &Sender{}
}

func (s *Sender) Send() {
	return
}
