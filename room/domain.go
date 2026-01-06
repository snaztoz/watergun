package room

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

func (d *domain) createRoom(id, name string) (*Model, error) {
	if id == "" {
		uuidV7, err := uuid.NewV7()
		if err != nil {
			log.Error("Failed to generate UUID", "err", err)
			return nil, err
		}

		id = uuidV7.String()
	}

	return d.store.createRoom(id, name)
}

func (d *domain) fetchRoom(id string) *Model {
	return d.store.fetchRoom(id)
}

func (d *domain) createParticipant(roomID, userID string, canPublish bool) (*ParticipantModel, error) {
	uuidV7, err := uuid.NewV7()
	if err != nil {
		log.Error("Failed to generate UUID", "err", err)
		return nil, err
	}

	return d.store.createParticipant(
		roomID,
		uuidV7.String(),
		userID,
		canPublish,
	)
}

func (d *domain) fetchParticipantsList(roomID string) []*ParticipantModel {
	return d.store.fetchParticipantsList(roomID)
}
