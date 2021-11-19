package translator

import "sync"

type jobs struct {
	sync.Mutex
	events map[string]chan struct{}
}

func newJobs() *jobs {
	return &jobs{
		events: make(map[string]chan struct{}),
	}
}

func (j *jobs) watch(confID string) <-chan struct{} {
	j.Lock()
	defer j.Unlock()

	trigger, ok := j.events[confID]
	if !ok {
		return nil
	}

	return trigger
}

func (j *jobs) create(confID string) {
	j.Lock()
	defer j.Unlock()

	_, ok := j.events[confID]
	if ok {
		return
	}

	j.events[confID] = make(chan struct{})
}

func (j *jobs) resolve(confID string) {
	j.Lock()
	defer j.Unlock()

	trigger, ok := j.events[confID]
	if !ok {
		return
	}

	trigger <- struct{}{}
	close(trigger)
	delete(j.events, confID)
}
