package room

import (
	"time"
)

type Model struct {
	ID           string              `json:"id"`
	Name         string              `json:"name"`
	Participants []*ParticipantModel `json:"participants"`
	CreatedAt    time.Time           `json:"created_at"`
	UpdatedAt    time.Time           `json:"updated_at"`
}

type ParticipantModel struct {
	ID         string    `json:"id"`
	UserID     string    `json:"user_id"`
	CanPublish bool      `json:"can_publish"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}
