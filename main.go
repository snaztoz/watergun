package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/snaztoz/watergun/socket"
)

const port = "8080"

func main() {
	server := &http.Server{
		Addr:    fmt.Sprintf(":%s", port),
		Handler: handler(),
	}

	go runServer(server)

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	<-quit
	stopServer(server)
}

func runServer(server *http.Server) {
	slog.Info(fmt.Sprintf("Server is listening at port %s", port))

	if err := server.ListenAndServe(); !errors.Is(err, http.ErrServerClosed) {
		slog.Error("Server failed to listen at the specified port", "err", err)
		return
	}
}

func stopServer(server *http.Server) {
	slog.Info("Shutting down the server...")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := server.Shutdown(shutdownCtx); err != nil {
		slog.Error("Failed to gracefully shutting down the server", "err", err)
		return
	}

	slog.Info("Server shutted down")
}

func handler() http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("/up", handleHealthCheck)
	mux.HandleFunc("/socket", handleWS())

	return mux
}

func handleHealthCheck(w http.ResponseWriter, r *http.Request) {
	if err := json.NewEncoder(w).Encode(map[string]any{
		"message": "up",
	}); err != nil {
		slog.Error("Failed to send message", "err", err)
	}
}

func handleWS() func(w http.ResponseWriter, r *http.Request) {
	hub := socket.NewHub()

	go hub.Run()

	return func(w http.ResponseWriter, r *http.Request) {
		socket.ServeWS(hub, w, r)
	}
}
