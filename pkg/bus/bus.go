package bus

import (
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
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

func (b Bus) Subscribe(event string) *Subscriber {
	b.lock.Lock()
	defer b.lock.Unlock()

	sub := &Subscriber{event: event, C: make(chan interface{}), clear: func() {
		b.clear(event)
	}, uid: uuid.NewString()}
	subs := b.subscribers[event]
	b.subscribers[event] = append(subs, sub)

	return sub
}

func (b Bus) Publish(event string, data interface{}) {
	b.lock.RLock()
	defer b.lock.RUnlock()

	subs := b.subscribers[event]

	logrus.Info("PUBLISH ", event)
	logrus.Infof("PUBLISH SUBS: %+v", subs)
	for _, sub := range subs {
		logrus.Info("PUBLISH sub", event, " ", sub.uid)
		if sub.C == nil {
			continue
		}
		sub.C <- data
	}
}

func (b *Bus) clear(event string) {
	logrus.Info("clear ", event)
	b.lock.Lock()
	defer b.lock.Unlock()

	subs := b.subscribers[event]
	logrus.Infof("SUBS: %+v",  subs)
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

	logrus.Infof("AFTER CLEAR SUBS: %+v",  b.subscribers[event])
}

var _ *Bus = &Bus{}
