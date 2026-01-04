package channel

import (
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestCreatingUserWithID(t *testing.T) {
	channelStore := NewStore()
	channelDomain := NewDomain(channelStore)

	id := "some-id-provided-externally"

	channel, err := channelDomain.createChannel(id, "some-name")
	assert.Nil(t, err)

	assert.Equal(t, id, channel.ID)
}

func TestCreatingUserWithoutID(t *testing.T) {
	userStore := NewStore()
	userDomain := NewDomain(userStore)

	channel, err := userDomain.createChannel("", "some-name")
	assert.Nil(t, err)

	assert.Nil(t, uuid.Validate(channel.ID))
}
