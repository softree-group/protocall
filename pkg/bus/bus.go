package bus

import (
	"sync"

	"github.com/hashicorp/go-uuid"
)

type Client interface {
	Subscribe(event string) *Subscriber
	Publish(event string, data interface{})
}

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

func (b Bus) Subscribe(event string) *Subscriber {
	b.lock.Lock()
	defer b.lock.Unlock()

	uuid, _ := uuid.GenerateUUID()
	sub := &Subscriber{
		event: event,
		C:     make(chan interface{}),
		clear: func() {
			b.clear(event)
		},
		uid: uuid,
	}
	subs := b.subscribers[event]
	b.subscribers[event] = append(subs, sub)

	return sub
}

func (b Bus) Publish(event string, data interface{}) {
	b.lock.RLock()
	defer b.lock.RUnlock()

	subs := b.subscribers[event]
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

	subs := b.subscribers[event]
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
