package applications

import "protocall/internal/connector/domain/entity"

type Socket interface {
	PublishConnectionEvent(user *entity.User) error
	PublishConnectedEvent(user *entity.User) error
	PublishLeaveEvent(user *entity.User) error
	PublishStartRecordEvent(conferenceID string) error
	PublishEndConference(conferenceID string) error
	PublishUserMessage(user *entity.User, message entity.SocketMessage) error
}
