package services

import "protocall/internal/connector/domain/entity"

type Socket interface {
	Publish(channel string, data entity.SocketMessage) error
}
