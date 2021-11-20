package translator

import (
	"fmt"
	"sync"
)

type jobs struct {
	sync.Mutex
	events map[string]chan struct{}
}

func newJobs() *jobs {
	return &jobs{
		events: make(map[string]chan struct{}),
	}
}

func (j *jobs) watch(confID []string) <-chan struct{} {
	j.Lock()
	defer j.Unlock()

	trigger, ok := j.events[confID]
	if !ok {
		return nil
	}

	return trigger
}

func (j *jobs) create(record string) {
	j.Lock()
	defer j.Unlock()

	_, ok := j.events[record]
	if ok {
		fmt.Println("found")
		return
	}

	j.events[record] = make(chan struct{})
}

func (j *jobs) resolve(record string) {
	j.Lock()
	defer j.Unlock()

	trigger, ok := j.events[record]
	if !ok {
		return
	}

	trigger <- struct{}{}
	close(trigger)
	delete(j.events, record)
}
