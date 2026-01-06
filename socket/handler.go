package socket

import (
	"net/http"

	"github.com/snaztoz/watergun/log"
)

func NewHandler(hub *hub) *handler {
	return &handler{hub: hub}
}

type handler struct {
	hub *hub
}

func (h *handler) Handle(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Error("Failed to upgrade to WebSocket", "err", err)
		return
	}

	client := newClient(h.hub, conn)

	h.hub.register <- client

	go client.pumpWrite()
	go client.pumpRead()
}
