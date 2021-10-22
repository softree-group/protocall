package memory

import (
	"protocall/domain/entity"
	"protocall/domain/repository"

	"github.com/google/btree"
)

type ConferenceMemory struct {
	store *btree.BTree
}

func (c ConferenceMemory) GetConference(conferenceID string) *entity.Conference {
	item := c.store.Get(&entity.Conference{ID: conferenceID})
	if item == nil {
		return nil
	}
	return item.(*entity.Conference)
}

func (c ConferenceMemory) SaveConference(conference *entity.Conference) {
	c.store.ReplaceOrInsert(conference)
}

func (c ConferenceMemory) DeleteConference(conferenceID string) {
	c.store.Delete(&entity.Conference{ID: conferenceID})
}

func NewConference() *ConferenceMemory {
	return &ConferenceMemory{store: btree.New(32)}
}

var _ repository.Conference = &ConferenceMemory{}
