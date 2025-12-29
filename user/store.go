package user

import (
	"errors"
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

func (s *store) createUser(id, masterID string) (*User, error) {
	user := &User{
		ID:        id,
		MasterID:  masterID,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	s.users[masterID] = user

	return user, nil
}

func (s *store) retrieveUser(masterID string) *User {
	user, exist := s.users[masterID]
	if !exist {
		return nil
	}

	return user
}

func NewValidator(store *store) *validator {
	return &validator{
		store: store,
	}
}

type validator struct {
	store *store
}

func (v *validator) validateUserCreation(masterID string) error {
	for _, user := range v.store.users {
		if user.MasterID == masterID {
			return errors.New("user already exist")
		}
	}
	return nil
}
