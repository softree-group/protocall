package services

import "protocall/pkg/bus"

type Bus interface {
	Subscribe(event string) *bus.Subscriber
	Publish(event string, data interface{})
}
