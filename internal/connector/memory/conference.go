package memory

import (
	"protocall/domain/entity"
	"protocall/domain/repository"
	"protocall/internal/config"

	"github.com/google/btree"
	"github.com/spf13/viper"
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
	return &ConferenceMemory{store: btree.New(viper.GetInt(config.Participant))}
}

var _ repository.Conference = &ConferenceMemory{}
