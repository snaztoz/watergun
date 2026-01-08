package user

import (
	"github.com/google/uuid"

	"github.com/snaztoz/watergun/log"
)

func NewDomain(store *store) *Domain {
	return &Domain{
		store: store,
	}
}

type Domain struct {
	store *store
}

func (d *Domain) createUser(id string) (*Model, error) {
	if id == "" {
		uuidV7, err := uuid.NewV7()
		if err != nil {
			log.Error("failed to generate UUID", "err", err)
			return nil, err
		}

		id = uuidV7.String()
	}

	return d.store.createUser(id)
}

func (d *Domain) FetchUser(id string) *Model {
	return d.store.fetchUser(id)
}
