package watergun

import "testing"

func TestSendingMessage(t *testing.T) {
	userID := "fake-user-id"
	channelID := "fake-channel-id"
	content := "some test message"

	d := NewMessageDomain(&dummyBroadcaster{}, &dummyMessageStorer{})

	message, err := d.SendMessage(userID, channelID, content)
	if err != nil {
		t.Fatalf("failed to send message: %v\n", err)
	}

	if message.UserID != userID || message.ChannelID != channelID || message.Content != content {
		t.Fatalf("attribute(s) mismatch: %v\n", message)
	}
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
