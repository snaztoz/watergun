package socket

import "testing"

func TestBlockingReadMessage(t *testing.T) {
	t.Run("unauthenticated & non-authenticating message", func(t *testing.T) {
		m := ReadMessage{}
		if !m.shouldBeBlocked() {
			t.Errorf("message should be blocked: %v\n", m)
		}
	})

	t.Run("authenticated & authenticating message", func(t *testing.T) {
		m := ReadMessage{
			UserID:      "some-id",
			MessageType: "auth",
		}

		if !m.shouldBeBlocked() {
			t.Errorf("message should be blocked: %v\n", m)
		}
	})
}
