package watergun

import "time"

type User struct {
	ID        string
	CreatedAt time.Time
	UpdatedAt *time.Time
}

func NewUserDomain(userCreator UserCreator) *UserDomain {
	return &UserDomain{userCreator: userCreator}
}

type UserDomain struct {
	userCreator UserCreator
}

func (d *UserDomain) CreateUser(id string) (*User, error) {
	return d.userCreator.CreateUser(id)
}

type UserCreator interface {
	CreateUser(id string) (*User, error)
}
