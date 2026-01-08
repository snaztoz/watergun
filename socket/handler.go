package socket

import (
	"net/http"

	"github.com/snaztoz/watergun/log"
	"github.com/snaztoz/watergun/response"
	"github.com/snaztoz/watergun/serverctx"
	"github.com/snaztoz/watergun/user"
)

func NewHandler(hub *hub, userDomain *user.Domain) *handler {
	return &handler{hub: hub, userDomain: userDomain}
}

type handler struct {
	hub *hub

	userDomain *user.Domain
}

func (h *handler) Handle(w http.ResponseWriter, r *http.Request) {
	id := r.Context().Value(serverctx.UserIDKey).(string)

	if user := h.userDomain.FetchUser(id); user == nil {
		log.Error("user does not exist", "id", id)
		response.SendErrorJSON(w, "user does not exist", 403)
		return
	}

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Error("failed to upgrade to WebSocket", "err", err)
		return
	}

	client := newClient(userID(id), h.hub, conn)

	h.hub.register <- client

	go client.pumpWrite()
	go client.pumpRead()
}
