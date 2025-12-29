package channel

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

func (d *domain) createChannel(id, name string) (*Channel, error) {
	if id == "" {
		uuidV7, err := uuid.NewV7()
		if err != nil {
			watergun.Logger().Error("Failed to generate UUID", "err", err)
			return nil, err
		}

		id = uuidV7.String()
	}

	return d.store.createChannel(id, name)
}

func (d *domain) retrieveChannel(masterID string) *Channel {
	return d.store.retrieveChannel(masterID)
}
