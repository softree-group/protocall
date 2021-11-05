package services

type Socket interface {
	Publish(channel string, data []byte) error
}
