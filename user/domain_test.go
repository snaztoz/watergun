package user

import (
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestCreatingUserWithID(t *testing.T) {
	userStore := NewStore()
	userDomain := NewDomain(userStore)

	id := "some-id-provided-externally"

	user, err := userDomain.createUser(id)
	assert.Nil(t, err)

	assert.Equal(t, id, user.ID)
}

func TestCreatingUserWithoutID(t *testing.T) {
	userStore := NewStore()
	userDomain := NewDomain(userStore)

	user, err := userDomain.createUser("")
	assert.Nil(t, err)

	assert.Nil(t, uuid.Validate(user.ID))
}
