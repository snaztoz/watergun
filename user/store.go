package user

import (
	"time"
)

func NewStore() *store {
	return &store{
		users: make(map[string]*User),
	}
}

type store struct {
	users map[string]*User
}

func (s *store) createUser(id string) (*User, error) {
	now := time.Now()
	user := &User{
		ID:        id,
		CreatedAt: now,
		UpdatedAt: now,
	}

	s.users[id] = user

	return user, nil
}

func (s *store) fetchUser(id string) *User {
	user, exist := s.users[id]
	if !exist {
		return nil
	}

	return user
}
