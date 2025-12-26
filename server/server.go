package server

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"strings"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/httplog/v3"

	"github.com/snaztoz/watergun"
	"github.com/snaztoz/watergun/socket"
)

func New(port string) *Server {
	return &Server{
		port: port,
		httpServer: &http.Server{
			Addr:    fmt.Sprintf(":%s", port),
			Handler: handler(),
		},
	}
}

type Server struct {
	httpServer *http.Server
	port       string
}

func (s *Server) Run() {
	watergun.Logger().Info(fmt.Sprintf("Server is listening at port %s", s.port))

	if err := s.httpServer.ListenAndServe(); !errors.Is(err, http.ErrServerClosed) {
		watergun.Logger().Error(
			"Server failed to listen at the specified port",
			"err", err,
		)
		return
	}
}

func (s *Server) Stop() {
	watergun.Logger().Info("Shutting down the server...")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := s.httpServer.Shutdown(shutdownCtx); err != nil {
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

	bindAdminRoutes(r)

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

func bindAdminRoutes(r *chi.Mux) {
	r.Route("/admin", func(r chi.Router) {
		r.Use(adminRoutesAuth)

		r.Route("/users", func(r chi.Router) {
			r.Get("/hello", func(w http.ResponseWriter, r *http.Request) {
				if err := json.NewEncoder(w).Encode(map[string]string{
					"user": "greeter",
				}); err != nil {
					http.Error(w, "something went wrong", 500)
					return
				}
			})
		})
	})
}

func adminRoutesAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authorization := r.Header.Get("Authorization")

		if !strings.HasPrefix(authorization, "Bearer ") {
			http.Error(w, http.StatusText(403), 403)
			return
		}

		if key := strings.TrimPrefix(
			authorization,
			"Bearer ",
		); key != watergun.AdminAPIKey() {
			http.Error(w, http.StatusText(403), 403)
			return
		}

		next.ServeHTTP(w, r)
	})
}
