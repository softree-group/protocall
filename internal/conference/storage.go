package conference

import (
	"github.com/gomodule/redigo/redis"
)

type Storage struct {
	*redis.Pool
}

func NewStorage(conn *redis.Pool) *Storage {
	return &Storage{conn}
}

func (c *Storage) Store(conferenceID, recName string) error {
	// conferenceJobs, _ := c.conferences.LoadOrStore(conferenceID, &sync.Map{})

	// conferenceJobs.(*sync.Map).Store(recName, false)
	return nil
}

func (c *Storage) DoneJob(conferenceID, recName string) error {
	// jobs, ok := c.conferences.Load(conferenceID)
	// if !ok {
	// 	return errors.New("no such conference")
	// }

	// jobs.(*sync.Map).Delete(recName)
	return nil
}

func (c *Storage) IsDone(conferenceID string) (bool, error) {
	// jobs, ok := c.conferences.Load(conferenceID)
	// if !ok {
	// 	return false, errors.New("no such conference")
	// }
	// count := 0
	// jobs.(*sync.Map).Range(func(key, value interface{}) bool {
	// 	count++
	// 	return false
	// })

	// if count == 0 {
	// 	c.conferences.Delete(conferenceID)
	// 	return true, nil
	// }
	return false, nil
}

func (s *Storage) GetConference(conferenceID string) *Conference {
	// item := c.store.Get(&Conference{ID: conferenceID})
	// if item == nil {
	// 	return nil
	// }
	return nil
}

func (s *Storage) SaveConference(conference *Conference) {
	// c.store.ReplaceOrInsert(conference)
}

func (s *Storage) DeleteConference(conferenceID string) {
	// c.store.Delete(&Conference{ID: conferenceID})
}

func (s *Storage) GetConferenceInfo(id string) (*ConferenceInfo, error) {
	// conferenceInfo := ConferenceInfo{
	// 	ID:           conference.ID,
	// 	HostID:       conference.HostUserID,
	// 	Participants: nil,
	// 	IsRecording:  conference.IsRecording,
	// 	StartedAt:    conference.Start.Unix(),
	// }

	// conference.Participants.Ascend(func(item btree.Item) bool {
	// 	if item == nil {
	// 		return false
	// 	}
	// 	user := item.(*user.User)
	// 	if user == nil {
	// 		return false
	// 	}
	// 	channel := ""
	// 	if user.Channel != nil {
	// 		channel = user.Channel.ID
	// 	}
	// 	participants = append(participants, UserInfo{
	// 		Name:    user.Username,
	// 		ID:      user.AsteriskAccount,
	// 		Channel: channel,
	// 	})
	// 	return true
	// })

	// conferenceInfo.Participants = participants
	return nil, nil
}
