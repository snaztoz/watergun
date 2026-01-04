package channel

import (
	"github.com/google/uuid"
	"github.com/snaztoz/watergun/log"
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
			log.Error("Failed to generate UUID", "err", err)
			return nil, err
		}

		id = uuidV7.String()
	}

	return d.store.createChannel(id, name)
}

func (d *domain) fetchChannel(id string) *Channel {
	return d.store.fetchChannel(id)
}

func (d *domain) createParticipant(channeldID, userID string, canPublish bool) (*Participant, error) {
	uuidV7, err := uuid.NewV7()
	if err != nil {
		log.Error("Failed to generate UUID", "err", err)
		return nil, err
	}

	return d.store.createParticipant(
		channeldID,
		uuidV7.String(),
		userID,
		canPublish,
	)
}

func (d *domain) fetchParticipantsList(channelID string) []*Participant {
	return d.store.fetchParticipantsList(channelID)
}
