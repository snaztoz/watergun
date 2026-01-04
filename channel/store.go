package channel

import (
	"time"
)

func NewStore() *store {
	return &store{
		channels: make(map[string]*Channel),
	}
}

type store struct {
	channels map[string]*Channel
}

func (s *store) createChannel(id, name string) (*Channel, error) {
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

func (s *store) fetchChannel(id string) *Channel {
	channel, exist := s.channels[id]
	if !exist {
		return nil
	}

	return channel
}

func (s *store) createParticipant(channelID, participantID, userID string, canPublish bool) (*Participant, error) {
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

func (s *store) fetchParticipantsList(channelID string) []*Participant {
	return s.channels[channelID].Participants
}
