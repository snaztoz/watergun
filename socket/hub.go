package socket

import (
	"encoding/json"

	"github.com/snaztoz/watergun/log"
	"github.com/snaztoz/watergun/room"
)

func NewHub(roomDomain *room.Domain) *hub {
	return &hub{
		roomDomain:     roomDomain,
		broadcastQueue: make(chan *broadcastMessage, 64),
		clients:        make(map[userID]*client),
		register:       make(chan *client),
		unregister:     make(chan *client),
	}
}

type hub struct {
	roomDomain *room.Domain

	broadcastQueue chan *broadcastMessage
	clients        map[userID]*client
	register       chan *client
	unregister     chan *client
}

func (h *hub) Run() {
	for {
		select {
		case client := <-h.register:
			h.registerClient(client)

		case client := <-h.unregister:
			h.unregisterClient(client)

		case msg := <-h.broadcastQueue:
			h.broadcast(msg)
		}
	}
}

func (h *hub) processMessage(senderID userID, msg *ReadMessage) {
	room := h.roomDomain.FetchRoom(msg.RoomID)
	if room == nil {
		log.Error("room does not exist", "senderID", senderID, "roomID", msg.RoomID)
		return
	}

	writeMsg := &WriteMessage{
		SenderID: string(senderID),
		RoomID:   msg.RoomID,
		Content:  msg.Content,
	}

	recipients := make([]userID, 0, len(room.Participants)-1)
	for _, participant := range room.Participants {
		if userID(participant.UserID) != senderID {
			recipients = append(recipients, userID(participant.UserID))
		}
	}

	h.broadcastQueue <- newBroadcastMessage(writeMsg, recipients)
}

func (h *hub) broadcast(msg *broadcastMessage) {
	for _, recipientID := range msg.recipients {
		c, exist := h.clients[recipientID]
		if !exist {
			return
		}

		select {
		case c.send <- msg.Bytes():
		default:
			close(c.send)
			delete(h.clients, c.userID)
		}
	}
}

func (h *hub) registerClient(c *client) {
	h.clients[c.userID] = c
}

func (h *hub) unregisterClient(c *client) {
	if _, exist := h.clients[c.userID]; !exist {
		return
	}

	delete(h.clients, c.userID)
	close(c.send)
}

func newBroadcastMessage(msg *WriteMessage, recipients []userID) *broadcastMessage {
	return &broadcastMessage{
		msg:        msg,
		recipients: recipients,
	}
}

type broadcastMessage struct {
	msg        *WriteMessage
	recipients []userID
	bytes      []byte
}

func (bm *broadcastMessage) Bytes() []byte {
	if len(bm.bytes) != 0 {
		return bm.bytes
	}

	bytes, err := json.Marshal(bm.msg)
	if err != nil {
		log.Error("failed to marshal message into JSON bytes", "err", err)
		return nil
	}

	bm.bytes = bytes

	return bm.bytes
}
