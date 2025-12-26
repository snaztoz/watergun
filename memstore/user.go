package memstore

import (
	"errors"
	"time"

	"github.com/snaztoz/watergun"
)

func NewUserStore() *userStore {
	return &userStore{
		users: make(map[string]*watergun.User),
	}
}

type userStore struct {
	users map[string]*watergun.User
}

func (s *userStore) ValidateUserCreation(masterID string) error {
	for _, user := range s.users {
		if user.MasterID == masterID {
			return errors.New("user already exist")
		}
	}
	return nil
}

func (s *userStore) CreateUser(id, masterID string) (*watergun.User, error) {
	user := &watergun.User{
		ID:        id,
		MasterID:  masterID,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	s.users[id] = user

	return user, nil
}
