package message

import (
	"time"

	"github.com/google/uuid"

	"github.com/snaztoz/watergun/log"
)

type Message struct {
	ID        string
	UserID    string
	RoomID    string
	Content   string
	CreatedAt time.Time
	UpdatedAt *time.Time
}

func NewMessageDomain(broadcaster Broadcaster, storer MessageStorer) *MessageDomain {
	return &MessageDomain{
		storer:      storer,
		broadcaster: broadcaster,
	}
}

type MessageDomain struct {
	storer      MessageStorer
	broadcaster Broadcaster
}

func (d *MessageDomain) SendMessage(userID, roomID, content string) (*Message, error) {
	id, err := uuid.NewV7()
	if err != nil {
		log.Error("failed to generate UUID", "err", err)
		return nil, err
	}

	m := &Message{
		ID:        id.String(),
		UserID:    userID,
		RoomID:    roomID,
		Content:   content,
		CreatedAt: time.Now(),
	}

	if err := d.storer.Store(m); err != nil {
		log.Error("failed to persist message", "err", err)
		return nil, err
	}

	if err := d.broadcaster.Broadcast(m); err != nil {
		log.Error("failed to broadcast message", "err", err)
		return nil, err
	}

	return m, nil
}

type MessageStorer interface {
	Store(message *Message) error
}

type Broadcaster interface {
	Broadcast(message *Message) error
}
