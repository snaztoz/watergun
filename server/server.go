package server

import (
	"context"
	"crypto"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/httplog/v3"

	"github.com/snaztoz/watergun/log"
	"github.com/snaztoz/watergun/room"
	"github.com/snaztoz/watergun/socket"
	"github.com/snaztoz/watergun/user"
)

const (
	defaultAdminKey = "ADMIN-INSECURE-KEY"
)

func New(port string, adminKey string, publicKey crypto.PublicKey) *Server {
	if adminKey == "" {
		adminKey = defaultAdminKey
	}

	return &Server{
		port:      port,
		adminKey:  adminKey,
		publicKey: publicKey,
	}
}

type Server struct {
	port      string
	adminKey  string
	publicKey crypto.PublicKey
	srv       *http.Server
}

func (s *Server) Run() {
	s.srv = &http.Server{
		Addr:    fmt.Sprintf(":%s", s.port),
		Handler: s.handler(),
	}

	log.Info("server is listening", "port", s.port)

	if err := s.srv.ListenAndServe(); !errors.Is(err, http.ErrServerClosed) {
		log.Error("server failed to listen at the specified port", "err", err)
		return
	}
}

func (s *Server) Stop() {
	log.Info("shutting down the server")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := s.srv.Shutdown(shutdownCtx); err != nil {
		log.Error("failed to gracefully shutting down the server", "err", err)
		return
	}

	log.Info("server shutted down")
}

func (s *Server) handler() http.Handler {
	r := chi.NewRouter()

	s.bootstrapMiddlewares(r)
	s.bootstrapRoutes(r)

	return r
}

func (*Server) bootstrapMiddlewares(r *chi.Mux) {
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)

	r.Use(httplog.RequestLogger(log.Logger(), &httplog.Options{
		Level:  slog.LevelInfo,
		Schema: httplog.SchemaECS,
	}))

	r.Use(middleware.Recoverer)
	r.Use(middleware.RedirectSlashes)
	r.Use(middleware.Timeout(60 * time.Second))

	r.Use(middleware.Heartbeat("/up"))
}

func (s *Server) bootstrapRoutes(r *chi.Mux) {
	roomStore := room.NewStore()
	roomDomain := room.NewDomain(roomStore)

	userStore := user.NewStore()
	userDomain := user.NewDomain(userStore)

	r.Route("/socket", func(r chi.Router) {
		r.Use(accessTokenParser(allowQueryParamToken))
		r.Use(socketRouteAuth(s.publicKey))

		socketHub := socket.NewHub(roomDomain)
		socketHandler := socket.NewHandler(socketHub, userDomain)

		go socketHub.Run()

		r.Get("/", socketHandler.Handle)
	})

	r.Route("/admin", func(r chi.Router) {
		r.Use(accessTokenParser(!allowQueryParamToken))
		r.Use(adminRoutesAuth(s.adminKey))
		r.Use(jsonContentType)

		r.Route("/rooms", func(r chi.Router) {
			roomHandler := room.NewHandler(roomDomain)

			r.Post("/", roomHandler.CreateRoom)
			r.Get("/{id}", roomHandler.FetchRoom)

			r.Route("/{roomID}/participants", func(r chi.Router) {
				r.Get("/", roomHandler.FetchParticipantsList)
				r.Post("/", roomHandler.CreateParticipant)
			})
		})

		r.Route("/users", func(r chi.Router) {
			userHandler := user.NewHandler(userDomain)

			r.Post("/", userHandler.CreateUser)
			r.Get("/{id}", userHandler.FetchUser)
		})
	})
}
