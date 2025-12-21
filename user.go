package watergun

import "time"

type User struct {
	ID         string
	IsOnline   bool
	LastSeenAt *time.Time
	CreatedAt  time.Time
	UpdatedAt  *time.Time
}
