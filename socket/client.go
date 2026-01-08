package socket

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"

	"github.com/snaztoz/watergun/log"
)

const (
	writeWait      = 10 * time.Second
	pongWait       = 60 * time.Second
	pingPeriod     = (pongWait * 9) / 10
	maxMessageSize = 1 << 12
)

var (
	newline = []byte{'\n'}
	space   = []byte{' '}

	upgrader = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
	}
)

type (
	clientID string
	userID   string
)

func newClient(userID userID, hub *hub, conn *websocket.Conn) *client {
	uuidV7, err := uuid.NewV7()
	if err != nil {
		log.Error("failed to generate client ID", "err", err)
		return nil
	}

	return &client{
		id:     clientID(uuidV7.String()),
		userID: userID,
		hub:    hub,
		conn:   conn,
		send:   make(chan []byte, 256),
	}
}

type client struct {
	id     clientID
	userID userID
	hub    *hub
	conn   *websocket.Conn
	send   chan []byte
}

func (c *client) pumpRead() {
	defer func() {
		c.hub.unregister <- c
		c.conn.Close()
	}()

	c.conn.SetReadLimit(maxMessageSize)
	c.conn.SetReadDeadline(time.Now().Add(pongWait))

	c.conn.SetPongHandler(func(string) error {
		c.conn.SetReadDeadline(time.Now().Add(pongWait))
		return nil
	})

	c.readMessages()
}

func (c *client) readMessages() {
	for {
		_, rawMsg, err := c.conn.ReadMessage()
		if err != nil {
			if isConnectionClosedUnexpectedly(err) {
				log.Error("connection closed unexpectedly", "err", err)
			}
			break
		}

		var msg ReadMessage
		if err := json.Unmarshal(rawMsg, &msg); err != nil {
			log.Error("failed to read message", "err", err)
			continue
		}

		c.hub.processMessage(c.userID, &msg)
	}
}

func (c *client) pumpWrite() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		c.conn.Close()
	}()

	for {
		select {
		case message, ok := <-c.send:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))

			if !ok {
				// The hub has closed the channel.
				c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			w, err := c.conn.NextWriter(websocket.TextMessage)
			if err != nil {
				return
			}
			w.Write(message)

			// Add queued chat messages to the current websocket message.
			n := len(c.send)
			for range n {
				w.Write(newline)
				w.Write(<-c.send)
			}

			if err := w.Close(); err != nil {
				return
			}

		case <-ticker.C:
			if err := c.pingConnection(); err != nil {
				return
			}
		}
	}
}

func (c *client) pingConnection() error {
	c.conn.SetWriteDeadline(time.Now().Add(writeWait))
	return c.conn.WriteMessage(websocket.PingMessage, nil)
}

func isConnectionClosedUnexpectedly(err error) bool {
	return websocket.IsUnexpectedCloseError(
		err,
		websocket.CloseGoingAway,
		websocket.CloseAbnormalClosure,
	)
}
