package app

import (
	"errors"
	"github.com/hashicorp/go-uuid"
	"protocall/application/applications"
	"protocall/domain/entity"
	"protocall/domain/repository"
)

type Conference struct {
	reps *repository.Repositories
}

func NewConference(reps *repository.Repositories) *Conference {
	return &Conference{reps: reps}
}

func (c Conference) StartConference(user *entity.User) (*entity.Conference, error) {
	id, _ := uuid.GenerateUUID()
	conference := &entity.Conference{
		ID:           id,
		Participants: []*entity.User{user},
		HostUserID:   user.AsteriskAccount,
		BridgeID:     "",
	}
	user.ConferenceID = id
	c.reps.User.Save(user)
	c.reps.Conference.Save(conference)
	return conference, nil
}

func (c Conference) JoinToConference(user *entity.User, meetID string) (*entity.Conference, error) {
	conference := c.reps.Conference.Get(meetID)
	if conference == nil {
		return nil, errors.New("no such meeting")
	}
	conference.Participants = append(conference.Participants, user)
	user.ConferenceID = meetID
	c.reps.User.Save(user)
	c.reps.Conference.Save(conference)
	return conference, nil
}

func (c Conference) IsExist(meetID string) bool {
	return c.reps.Conference.Get(meetID) != nil
}

var _ applications.Conference = Conference{}
