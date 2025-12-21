package watergun

import (
	"log/slog"
	"time"

	"github.com/google/uuid"
)

type Message struct {
	ID        string
	UserID    string
	ChannelID string
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

func (d *MessageDomain) SendMessage(userID, channelID, content string) (*Message, error) {
	id, err := uuid.NewV7()
	if err != nil {
		slog.Error("Failed to generate UUID", "err", err)
		return nil, err
	}

	m := &Message{
		ID:        id.String(),
		UserID:    userID,
		ChannelID: channelID,
		Content:   content,
		CreatedAt: time.Now(),
	}

	if err := d.storer.Store(m); err != nil {
		slog.Error("Failed to persist message", "err", err)
		return nil, err
	}

	if err := d.broadcaster.Broadcast(m); err != nil {
		slog.Error("Failed to broadcast message", "err", err)
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
