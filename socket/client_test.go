package socket

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBlockingReadMessage(t *testing.T) {
	t.Run("unauthenticated & non-authenticating message", func(t *testing.T) {
		m := ReadMessage{}

		assert.True(t, m.shouldBeBlocked())
	})

	t.Run("authenticated & authenticating message", func(t *testing.T) {
		m := ReadMessage{
			UserID:      "some-id",
			MessageType: "auth",
		}

		assert.True(t, m.shouldBeBlocked())
	})
}
