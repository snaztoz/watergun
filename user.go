package watergun

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID        string    `json:"-"`
	MasterID  string    `json:"master_id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func NewUserDomain(storer userStorer, validator userValidator) *UserDomain {
	return &UserDomain{
		storer:    storer,
		validator: validator,
	}
}

type UserDomain struct {
	storer    userStorer
	validator userValidator
}

func (d *UserDomain) CreateUser(masterID string) (*User, error) {
	if err := d.validator.ValidateUserCreation(masterID); err != nil {
		logger.Error("Validation failed", "err", err)
		return nil, err
	}

	id, err := uuid.NewV7()
	if err != nil {
		logger.Error("Failed to generate UUID", "err", err)
		return nil, err
	}

	return d.storer.CreateUser(id.String(), masterID)
}

func (d *UserDomain) RetrieveUser(masterID string) *User {
	return d.storer.RetrieveUser(masterID)
}

type userStorer interface {
	CreateUser(id, masterID string) (*User, error)
	RetrieveUser(masterID string) *User
}

type userValidator interface {
	ValidateUserCreation(masterID string) error
}
