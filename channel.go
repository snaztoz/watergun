package watergun

import (
	"log/slog"
	"time"
)

type Channel struct {
	ID           string
	Name         string
	Participants map[string]*Participant
	CreatedAt    time.Time
	UpdatedAt    *time.Time
}

type Participant struct {
	ID        string
	UserID    string
	CreatedAt time.Time
	UpdatedAt *time.Time
}

func NewChannelDomain(
	channelCreator ChannelCreator,
	channelParticipantAdder ChannelParticipantAdder,
) *ChannelDomain {
	return &ChannelDomain{
		channelCreator:          channelCreator,
		channelParticipantAdder: channelParticipantAdder,
	}
}

type ChannelDomain struct {
	channelCreator          ChannelCreator
	channelParticipantAdder ChannelParticipantAdder
}

func (d *ChannelDomain) CreateChannel(name string) (*Channel, error) {
	if err := d.channelCreator.ValidateChannelCreation(name); err != nil {
		slog.Error("Validation failed", "err", err)
		return nil, err
	}

	return d.channelCreator.CreateChannel(name)
}

func (d *ChannelDomain) AddParticipant(userID string) (*Channel, error) {
	if err := d.channelParticipantAdder.ValidateParticipant(userID); err != nil {
		slog.Error("Validation failed", "err", err)
		return nil, err
	}

	return d.channelParticipantAdder.AddParticipant(userID)
}

type ChannelCreator interface {
	ValidateChannelCreation(name string) error
	CreateChannel(name string) (*Channel, error)
}

type ChannelParticipantAdder interface {
	ValidateParticipant(userID string) error
	AddParticipant(userID string) (*Channel, error)
}
