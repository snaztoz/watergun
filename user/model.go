package user

import (
	"time"
)

type User struct {
	ID        string    `json:"-"`
	MasterID  string    `json:"master_id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
