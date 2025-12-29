package user

import (
	"github.com/google/uuid"

	"github.com/snaztoz/watergun"
)

func NewDomain(store *store, validator *validator) *domain {
	return &domain{
		store:     store,
		validator: validator,
	}
}

type domain struct {
	store     *store
	validator *validator
}

func (d *domain) createUser(masterID string) (*User, error) {
	if err := d.validator.validateUserCreation(masterID); err != nil {
		watergun.Logger().Error("Validation failed", "err", err)
		return nil, err
	}

	id, err := uuid.NewV7()
	if err != nil {
		watergun.Logger().Error("Failed to generate UUID", "err", err)
		return nil, err
	}

	return d.store.createUser(id.String(), masterID)
}

func (d *domain) retrieveUser(masterID string) *User {
	return d.store.retrieveUser(masterID)
}
