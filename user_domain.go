package watergun

import (
	"github.com/google/uuid"
)

func NewUserDomain(store *userStore) *UserDomain {
	return &UserDomain{
		store: store,
	}
}

type UserDomain struct {
	store *userStore
}

func (d *UserDomain) CreateUser(id string) (*User, error) {
	if id == "" {
		uuidV7, err := uuid.NewV7()
		if err != nil {
			logger.Error("Failed to generate UUID", "err", err)
			return nil, err
		}

		id = uuidV7.String()
	}

	return d.store.createUser(id)
}

func (d *UserDomain) RetrieveUser(masterID string) *User {
	return d.store.retrieveUser(masterID)
}
