package app

import (
	"protocall/internal/connector/application/applications"
	"protocall/internal/connector/domain/entity"
	"protocall/internal/connector/domain/services"
)

type Socket struct {
	socketService services.Socket
}

func NewSocket(socketService services.Socket) *Socket {
	return &Socket{socketService}
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

func (s Socket) PublishUserMessage(user *entity.User, message entity.SocketMessage) error {
	return s.socketService.Publish("notify#"+user.AsteriskAccount, message)
}

func (s Socket) publishConferenceEventWithUserData(event string, user *entity.User) error {
	data, err := s.userMessage(event, user)
	if err != nil {
		return err
	}

	return s.socketService.Publish("conference~"+user.ConferenceID, data)
}

func (s Socket) publishConferenceEvent(event, conferenceID string) error {
	return s.socketService.Publish("conference~"+conferenceID, entity.SocketMessage{
		"event": event,
	})
}

func (s Socket) userMessage(event string, user *entity.User) (entity.SocketMessage, error) {
	var userChannel string
	if user.Channel != nil {
		userChannel = user.Channel.ID
	}
	payload := entity.SocketMessage{
		"event": event,
		"user": entity.SocketMessage{
			"id":      user.AsteriskAccount,
			"name":    user.Username,
			"channel": userChannel,
		},
	}
	return payload, nil
}

var _ applications.Socket = &Socket{}
