package user

import (
	"github.com/gomodule/redigo/redis"
)

type Storage struct {
	*redis.Pool
}

func NewStorage(conn *redis.Pool) *Storage {
	return &Storage{conn}
}

const pattern = "user:"

func (s *Storage) FindUser(sessionID string) *User {
	return nil
}

func (s *Storage) SaveUser(user *User) {

}

func (s *Storage) DeleteUser(sessionID string) {
}
