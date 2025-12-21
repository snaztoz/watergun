package watergun

import "time"

type Channel struct {
	ID           string
	Participants map[string]*Participant
	IsGroup      bool
	CreatedAt    time.Time
	UpdatedAt    *time.Time
}

type Participant struct {
	ID        string
	UserID    string
	Role      string
	CreatedAt time.Time
	UpdatedAt *time.Time
}
