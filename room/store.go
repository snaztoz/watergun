package room

import (
	"time"
)

func NewStore() *store {
	return &store{
		rooms: make(map[string]*Model),
	}
}

type store struct {
	rooms map[string]*Model
}

func (s *store) createRoom(id, name string) (*Model, error) {
	now := time.Now()
	room := &Model{
		ID:           id,
		Name:         name,
		Participants: make([]*ParticipantModel, 0, 2),
		CreatedAt:    now,
		UpdatedAt:    now,
	}

	s.rooms[id] = room

	return room, nil
}

func (s *store) fetchRoom(id string) *Model {
	room, exist := s.rooms[id]
	if !exist {
		return nil
	}

	return room
}

func (s *store) createParticipant(roomID, participantID, userID string, canPublish bool) (*ParticipantModel, error) {
	for _, participant := range s.rooms[roomID].Participants {
		if participant.UserID == userID {
			// no-op
			return participant, nil
		}
	}

	now := time.Now()
	participant := &ParticipantModel{
		ID:         participantID,
		UserID:     userID,
		CanPublish: canPublish,
		CreatedAt:  now,
		UpdatedAt:  now,
	}

	s.rooms[roomID].Participants = append(
		s.rooms[roomID].Participants,
		participant,
	)

	return participant, nil
}

func (s *store) fetchParticipantsList(roomID string) []*ParticipantModel {
	return s.rooms[roomID].Participants
}
