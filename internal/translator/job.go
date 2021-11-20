package translator

import (
	"sync"
)

type job struct {
	sync.Mutex
	events map[string]chan struct{}
}

func newJob() *job {
	return &job{
		events: make(map[string]chan struct{}),
	}
}

func (j *job) wait(ids []string) <-chan struct{} {
	j.Lock()
	defer j.Unlock()

	out := make(chan struct{})

	var wg sync.WaitGroup
	for _, id := range ids {
		trigger, ok := j.events[id]
		if !ok {
			continue
		}

		wg.Add(1)
		go func(c <-chan struct{}) {
			defer wg.Done()
			for range c {
			}
		}(trigger)
	}

	go func() {
		wg.Wait()
		out <- struct{}{}
		close(out)
	}()
	return out
}

func (j *job) create(id string) {
	j.Lock()
	defer j.Unlock()

	_, ok := j.events[id]
	if ok {
		return
	}

	j.events[id] = make(chan struct{})
}

func (j *job) resolve(id string) {
	j.Lock()
	defer j.Unlock()

	trigger, ok := j.events[id]
	if !ok {
		return
	}

	trigger <- struct{}{}
	close(trigger)
	delete(j.events, id)
}
