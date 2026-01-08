package user

import (
	"time"
)

func NewStore() *store {
	return &store{
		users: make(map[string]*Model),
	}
}

type store struct {
	users map[string]*Model
}

func (s *store) createUser(id string) (*Model, error) {
	now := time.Now()
	user := &Model{
		ID:        id,
		CreatedAt: now,
		UpdatedAt: now,
	}

	s.users[id] = user

	return user, nil
}

func (s *store) fetchUser(id string) *Model {
	user, exist := s.users[id]
	if !exist {
		return nil
	}

	return user
}
