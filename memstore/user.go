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

func (s *userStore) CreateUser(id, masterID string) (*watergun.User, error) {
	user := &watergun.User{
		ID:        id,
		MasterID:  masterID,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	s.users[masterID] = user

	return user, nil
}

func (s *userStore) RetrieveUser(masterID string) *watergun.User {
	user, exist := s.users[masterID]
	if !exist {
		return nil
	}

	return user
}

func NewUserValidator(store userStore) *userValidator {
	return &userValidator{
		store: store,
	}
}

type userValidator struct {
	store userStore
}

func (v *userValidator) ValidateUserCreation(masterID string) error {
	for _, user := range v.store.users {
		if user.MasterID == masterID {
			return errors.New("user already exist")
		}
	}
	return nil
}
