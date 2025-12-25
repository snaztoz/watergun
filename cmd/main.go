package main

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/httplog/v3"

	"github.com/snaztoz/watergun"
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
	watergun.Logger().Info(fmt.Sprintf("Server is listening at port %s", port))

	if err := server.ListenAndServe(); !errors.Is(err, http.ErrServerClosed) {
		watergun.Logger().Error(
			"Server failed to listen at the specified port",
			"err", err,
		)
		return
	}
}

func stopServer(server *http.Server) {
	watergun.Logger().Info("Shutting down the server...")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := server.Shutdown(shutdownCtx); err != nil {
		watergun.Logger().Error(
			"Failed to gracefully shutting down the server",
			"err", err,
		)
		return
	}

	watergun.Logger().Info("Server shutted down")
}

func handler() http.Handler {
	r := chi.NewRouter()

	bindMiddlewares(r)

	r.Get("/socket", handleWS())

	return r
}

func handleWS() func(w http.ResponseWriter, r *http.Request) {
	hub := socket.NewHub()

	go hub.Run()

	return func(w http.ResponseWriter, r *http.Request) {
		socket.ServeWS(hub, w, r)
	}
}

func bindMiddlewares(r *chi.Mux) {
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)

	r.Use(httplog.RequestLogger(watergun.Logger(), &httplog.Options{
		Level:  slog.LevelInfo,
		Schema: httplog.SchemaECS,
	}))

	r.Use(middleware.Recoverer)
	r.Use(middleware.RedirectSlashes)
	r.Use(middleware.Timeout(60 * time.Second))

	r.Use(middleware.Heartbeat("/up"))
}
