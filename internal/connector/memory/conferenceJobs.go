package memory

import (
	"errors"
	"sync"
)

type ConferenceJobs struct {
	conferences sync.Map
}

func NewConferenceJobs() *ConferenceJobs {
	return &ConferenceJobs{conferences: sync.Map{}}
}

func (c *ConferenceJobs) Store(conferenceID, recName string) error {
	conferenceJobs, _ := c.conferences.LoadOrStore(conferenceID, &sync.Map{})

	conferenceJobs.(*sync.Map).Store(recName, false)
	return nil
}

func (c *ConferenceJobs) DoneJob(conferenceID, recName string) error {
	jobs, ok := c.conferences.Load(conferenceID)
	if !ok {
		return errors.New("no such conference")
	}

	jobs.(*sync.Map).Delete(recName)
	return nil
}

func (c *ConferenceJobs) IsDone(conferenceID string) (bool, error) {
	jobs, ok := c.conferences.Load(conferenceID)
	if !ok {
		return false, errors.New("no such conference")
	}
	count := 0
	jobs.(*sync.Map).Range(func(key, value interface{}) bool {
		count++
		return false
	})

	if count == 0 {
		c.conferences.Delete(conferenceID)
		return true, nil
	}
	return false, nil
}
