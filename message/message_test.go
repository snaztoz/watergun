package message

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSendingMessage(t *testing.T) {
	userID := "fake-user-id"
	channelID := "fake-channel-id"
	content := "some test message"

	d := NewMessageDomain(&dummyBroadcaster{}, &dummyMessageStorer{})

	message, err := d.SendMessage(userID, channelID, content)
	assert.Nil(t, err)

	assert.Equal(t, userID, message.UserID)
	assert.Equal(t, channelID, message.ChannelID)
	assert.Equal(t, content, message.Content)
}

type dummyMessageStorer struct{}

func (*dummyMessageStorer) Store(m *Message) error {
	// no-op
	return nil
}

type dummyBroadcaster struct{}

func (*dummyBroadcaster) Broadcast(m *Message) error {
	// no-op
	return nil
}
