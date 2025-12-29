package user

import (
	"github.com/google/uuid"

	"github.com/snaztoz/watergun"
)

func NewDomain(store *store) *domain {
	return &domain{
		store: store,
	}
}

type domain struct {
	store *store
}

func (d *domain) createUser(id string) (*User, error) {
	if id == "" {
		uuidV7, err := uuid.NewV7()
		if err != nil {
			watergun.Logger().Error("Failed to generate UUID", "err", err)
			return nil, err
		}

		id = uuidV7.String()
	}

	return d.store.createUser(id)
}

func (d *domain) retrieveUser(masterID string) *User {
	return d.store.retrieveUser(masterID)
}
