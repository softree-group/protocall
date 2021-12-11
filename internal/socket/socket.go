package socket

import (
	"protocall/internal/user"
	"protocall/pkg/centrifugo"
)

const (
	SocketEventLeave       = "leave"
	SocketEventConnection  = "connection"
	SocketEventConnected   = "connected"
	SocketEventStartRecord = "start_record"
)

type Message = centrifugo.SocketMessage

type provider interface {
	Publish(channel string, data Message) error
}

type Socket struct {
	provider
}

func NewSocket(socketService provider) *Socket {
	return &Socket{socketService}
}

func (s Socket) publishConferenceEventWithUserData(event string, client *user.User) error {
	var userChannel string
	if client.Channel != nil {
		userChannel = client.Channel.ID
	}
	return s.Publish("conference:"+client.ConferenceID, Message{
		"event": event,
		"user": Message{
			"id":      client.AsteriskAccount,
			"name":    client.Username,
			"channel": userChannel,
		},
	})
}

func (s Socket) PublishConnectionEvent(client *user.User) error {
	return s.publishConferenceEventWithUserData("connection", client)
}

func (s Socket) PublishConnectedEvent(client *user.User) error {
	return s.publishConferenceEventWithUserData("connected", client)
}

func (s Socket) PublishLeaveEvent(client *user.User) error {
	return s.publishConferenceEventWithUserData("leave", client)
}

func (s Socket) publishConferenceEvent(event, conferenceID string) error {
	return s.Publish("conference:"+conferenceID, Message{
		"event": event,
	})
}

func (s Socket) PublishStartRecordEvent(conferenceID string) error {
	return s.publishConferenceEvent("start_record", conferenceID)
}

func (s Socket) PublishEndConference(conferenceID string) error {
	return s.publishConferenceEvent("end", conferenceID)
}

func (s Socket) PublishUserMessage(client *user.User, message Message) error {
	return s.Publish("notify#"+client.AsteriskAccount, message)
}
