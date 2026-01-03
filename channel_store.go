package watergun

import (
	"time"
)

func NewChannelStore() *channelStore {
	return &channelStore{
		channels: make(map[string]*Channel),
	}
}

type channelStore struct {
	channels map[string]*Channel
}

func (s *channelStore) createChannel(id, name string) (*Channel, error) {
	now := time.Now()
	channel := &Channel{
		ID:           id,
		Name:         name,
		Participants: make([]*Participant, 0, 2),
		CreatedAt:    now,
		UpdatedAt:    now,
	}

	s.channels[id] = channel

	return channel, nil
}

func (s *channelStore) fetchChannel(id string) *Channel {
	channel, exist := s.channels[id]
	if !exist {
		return nil
	}

	return channel
}

func (s *channelStore) createParticipant(channelID, participantID, userID string, canPublish bool) (*Participant, error) {
	for _, participant := range s.channels[channelID].Participants {
		if participant.UserID == userID {
			// no-op
			return participant, nil
		}
	}

	now := time.Now()
	participant := &Participant{
		ID:         participantID,
		UserID:     userID,
		CanPublish: canPublish,
		CreatedAt:  now,
		UpdatedAt:  now,
	}

	s.channels[channelID].Participants = append(
		s.channels[channelID].Participants,
		participant,
	)

	return participant, nil
}

func (s *channelStore) fetchParticipantsList(channelID string) []*Participant {
	return s.channels[channelID].Participants
}
