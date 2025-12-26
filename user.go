package watergun

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID        string    `json:"id"`
	MasterID  string    `json:"master_id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func NewUserDomain(userCreator UserCreator) *UserDomain {
	return &UserDomain{userCreator: userCreator}
}

type UserDomain struct {
	userCreator UserCreator
}

func (d *UserDomain) CreateUser(masterID string) (*User, error) {
	if err := d.userCreator.ValidateUserCreation(masterID); err != nil {
		logger.Error("Validation failed", "err", err)
		return nil, err
	}

	id, err := uuid.NewV7()
	if err != nil {
		logger.Error("Failed to generate UUID", "err", err)
		return nil, err
	}

	return d.userCreator.CreateUser(id.String(), masterID)
}

type UserCreator interface {
	ValidateUserCreation(masterID string) error
	CreateUser(id, masterID string) (*User, error)
}
