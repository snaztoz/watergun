package socket

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
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
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func ServeWS(hub *Hub, w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		slog.Error("Failed to upgrade to WebSocket", "err", err)
		return
	}

	client := &Client{
		hub:  hub,
		conn: conn,
		send: make(chan []byte, 256),
	}

	client.hub.register <- client

	go client.setWritePump()
	go client.setReadPump()
}

type Client struct {
	hub    *Hub
	conn   *websocket.Conn
	send   chan []byte
	userID string
}

func (c *Client) setReadPump() {
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

func (c *Client) readMessages() {
	for {
		_, rawMessage, err := c.conn.ReadMessage()
		if err != nil {
			if c.isConnectionClosedUnexpectedly(err) {
				slog.Error("Connection closed unexpectedly", "err", err)
			}
			break
		}

		var message ReadMessage
		if err := json.Unmarshal(rawMessage, &message); err != nil {
			slog.Error("Failed to read message", "err", err)
			continue
		}
		message.UserID = c.userID

		if message.shouldBeBlocked() {
			continue
		}

		c.hub.broadcast <- []byte(message.Content)
	}
}

func (c *Client) setWritePump() {
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
				// The hub closed the channel.
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
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))

			if err := c.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}

func (c *Client) isConnectionClosedUnexpectedly(err error) bool {
	return websocket.IsUnexpectedCloseError(
		err,
		websocket.CloseGoingAway,
		websocket.CloseAbnormalClosure,
	)
}

type ReadMessage struct {
	UserID      string `json:"-"`
	MessageType string `json:"type"`
	Content     string `json:"content"`
}

func (m *ReadMessage) isAuthenticated() bool {
	return m.UserID != ""
}

func (m *ReadMessage) isForAuthentication() bool {
	return m.MessageType == "auth"
}

func (m *ReadMessage) shouldBeBlocked() bool {
	return (!m.isAuthenticated() && !m.isForAuthentication()) ||
		(m.isAuthenticated() && m.isForAuthentication())
}
