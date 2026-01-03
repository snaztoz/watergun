package watergun

import (
	"time"
)

func NewUserStore() *userStore {
	return &userStore{
		users: make(map[string]*User),
	}
}

type userStore struct {
	users map[string]*User
}

func (s *userStore) createUser(id string) (*User, error) {
	now := time.Now()
	user := &User{
		ID:        id,
		CreatedAt: now,
		UpdatedAt: now,
	}

	s.users[id] = user

	return user, nil
}

func (s *userStore) fetchUser(id string) *User {
	user, exist := s.users[id]
	if !exist {
		return nil
	}

	return user
}
