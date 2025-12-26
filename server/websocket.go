package server

import (
	"net/http"

	"github.com/snaztoz/watergun/socket"
)

func handleWS() func(w http.ResponseWriter, r *http.Request) {
	hub := socket.NewHub()

	go hub.Run()

	return func(w http.ResponseWriter, r *http.Request) {
		socket.ServeWS(hub, w, r)
	}
}
