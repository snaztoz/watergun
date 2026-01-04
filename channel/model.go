package channel

import (
	"time"
)

type Channel struct {
	ID           string         `json:"id"`
	Name         string         `json:"name"`
	Participants []*Participant `json:"participants"`
	CreatedAt    time.Time      `json:"created_at"`
	UpdatedAt    time.Time      `json:"updated_at"`
}

type Participant struct {
	ID         string    `json:"id"`
	UserID     string    `json:"user_id"`
	CanPublish bool      `json:"can_publish"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}
