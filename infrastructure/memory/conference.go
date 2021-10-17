package memory

import (
	"github.com/google/btree"
	"protocall/domain/entity"
	"protocall/domain/repository"
)

type ConferenceMemory struct {
	store *btree.BTree
}

func (c ConferenceMemory) Get(conferenceID string) *entity.Conference {
	item := c.store.Get(&entity.Conference{ID: conferenceID})
	if item == nil {
		return nil
	}
	return item.(*entity.Conference)
}

func (c ConferenceMemory) Save(conference *entity.Conference) {
	c.store.ReplaceOrInsert(conference)
}

func (c ConferenceMemory) Delete(conferenceID string) {
	c.store.Delete(&entity.Conference{ID: conferenceID})
}

func NewConference() *ConferenceMemory {
	return &ConferenceMemory{store: btree.New(32)}
}

var _ repository.Conference = &ConferenceMemory{}
