package app

import (
	"encoding/json"
	"protocall/application/applications"
	"protocall/domain/entity"
	"protocall/domain/services"
)

type Socket struct {
	socketService services.Socket
}

func (s Socket) PublishConnectionEvent(user *entity.User) error {
	return s.publishConferenceEventWithUserData("connection", user)
}

func (s Socket) PublishConnectedEvent(user *entity.User) error {
	return s.publishConferenceEventWithUserData("connected", user)
}

func (s Socket) PublishLeaveEvent(user *entity.User) error {
	return s.publishConferenceEventWithUserData("leave", user)
}

func (s Socket) PublishStartRecordEvent(conferenceID string) error {
	return s.publishConferenceEvent("start_record", conferenceID)
}

func (s Socket) PublishEndConference(conferenceID string) error {
	return s.publishConferenceEvent("end", conferenceID)
}

func (s Socket) publishConferenceEventWithUserData(event string, user *entity.User) error {
	data, err := s.userMessage(event, user)
	if err != nil {
		return err
	}

	return s.socketService.Publish("conference~" + user.ConferenceID, data)
}

func (s Socket) publishConferenceEvent(event, conferenceID string) error {
	data, err := s.eventMessage(event)
	if err != nil {
		return err
	}
	return s.socketService.Publish("conference~" + conferenceID, data)
}

func (s Socket) userMessage(event string, user *entity.User) ([]byte, error) {
	var userChannel string
	if user.Channel != nil {
		userChannel = user.Channel.ID
	}
	payload := map[string]interface{}{
		"event": event,
		"user": map[string]interface{}{
			"id": user.AsteriskAccount,
			"name": user.Username,
			"channel": userChannel,
		},
	}

	return json.Marshal(payload)
}

func (s Socket) eventMessage(event string) ([]byte, error) {
	return json.Marshal(map[string]interface{}{"event": event})
}

var _ applications.Socket = &Socket{}
