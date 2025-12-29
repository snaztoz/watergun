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

func (s *store) retrieveChannel(id string) *Channel {
	channel, exist := s.channels[id]
	if !exist {
		return nil
	}

	return channel
}
