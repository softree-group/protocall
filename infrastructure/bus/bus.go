package bus

import (
	"protocall/domain/services"
	"sync"
)

type Bus struct {
	subscribers map[string][]*Subscriber
	lock        *sync.RWMutex
}

func New() *Bus {
	return &Bus{
		subscribers: map[string][]*Subscriber{},
		lock:        &sync.RWMutex{},
	}
}

func (b Bus) Subscribe(event string) services.Subscriber {
	b.lock.Lock()
	defer b.lock.Unlock()

	sub := &Subscriber{C: make(chan interface{}), clear: func() {
		b.clear(event)
	}}
	subs, _ := b.subscribers[event]
	b.subscribers[event] = append(subs, sub)

	return sub
}

func (b Bus) Publish(event string, data interface{}) {
	b.lock.RLock()
	defer b.lock.RUnlock()

	subs, _ := b.subscribers[event]
	if len(subs) == 0 {
		return
	}

	for _, sub := range subs {
		if sub.C == nil {
			continue
		}
		sub.C <- data
	}
}

func (b *Bus) clear(event string) {
	b.lock.Lock()
	defer b.lock.Unlock()

	subs, _ := b.subscribers[event]
	if len(subs) == 0 {
		return
	}

	var notNilSubs []*Subscriber

	for _, sub := range subs {
		if sub.C == nil {
			continue
		}
		notNilSubs = append(notNilSubs, sub)
	}

	if len(notNilSubs) == 0 {
		delete(b.subscribers, event)
	}

	b.subscribers[event] = notNilSubs
}

var _ services.Bus = &Bus{}
