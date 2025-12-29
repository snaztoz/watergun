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

func (d *domain) create(id, name string) (*Channel, error) {
	if id == "" {
		uuidV7, err := uuid.NewV7()
		if err != nil {
			watergun.Logger().Error("Failed to generate UUID", "err", err)
			return nil, err
		}

		id = uuidV7.String()
	}

	return d.store.create(id, name)
}

func (d *domain) retrieve(id string) *Channel {
	return d.store.retrieve(id)
}

func (d *domain) createParticipant(
	channeldID string,
	userID string,
	canPublish bool,
) (*Participant, error) {
	uuidV7, err := uuid.NewV7()
	if err != nil {
		watergun.Logger().Error("Failed to generate UUID", "err", err)
		return nil, err
	}

	return d.store.createParticipant(
		channeldID,
		uuidV7.String(),
		userID,
		canPublish,
	)
}

func (d *domain) retrieveParticipantsList(channelID string) []*Participant {
	return d.store.retrieveParticipantsList(channelID)
}
