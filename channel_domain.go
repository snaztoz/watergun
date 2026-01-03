package watergun

import (
	"github.com/google/uuid"
)

func NewChannelDomain(store *channelStore) *ChannelDomain {
	return &ChannelDomain{
		store: store,
	}
}

type ChannelDomain struct {
	store *channelStore
}

func (d *ChannelDomain) CreateChannel(id, name string) (*Channel, error) {
	if id == "" {
		uuidV7, err := uuid.NewV7()
		if err != nil {
			logger.Error("Failed to generate UUID", "err", err)
			return nil, err
		}

		id = uuidV7.String()
	}

	return d.store.createChannel(id, name)
}

func (d *ChannelDomain) Fetch(id string) *Channel {
	return d.store.fetchChannel(id)
}

func (d *ChannelDomain) CreateParticipant(channeldID, userID string, canPublish bool) (*Participant, error) {
	uuidV7, err := uuid.NewV7()
	if err != nil {
		logger.Error("Failed to generate UUID", "err", err)
		return nil, err
	}

	return d.store.createParticipant(
		channeldID,
		uuidV7.String(),
		userID,
		canPublish,
	)
}

func (d *ChannelDomain) FetchParticipantsList(channelID string) []*Participant {
	return d.store.fetchParticipantsList(channelID)
}
