package room

import (
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestCreatingUserWithID(t *testing.T) {
	roomStore := NewStore()
	roomDomain := NewDomain(roomStore)

	id := "some-id-provided-externally"

	room, err := roomDomain.createRoom(id, "some-name")
	assert.Nil(t, err)

	assert.Equal(t, id, room.ID)
}

func TestCreatingUserWithoutID(t *testing.T) {
	userStore := NewStore()
	userDomain := NewDomain(userStore)

	room, err := userDomain.createRoom("", "some-name")
	assert.Nil(t, err)

	assert.Nil(t, uuid.Validate(room.ID))
}
