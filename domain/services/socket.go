package services

import "protocall/domain/entity"

type Socket interface {
	Publish(channel string, data entity.SocketMessage) error
}
