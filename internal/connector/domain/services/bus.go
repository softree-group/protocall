package services

type Subscriber interface {
	Cancel()
	Channel() chan interface{}
}

type Bus interface {
	Subscribe(event string) Subscriber
	Publish(event string, data interface{})
}
