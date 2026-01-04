package server

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"strings"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/httplog/v3"

	"github.com/snaztoz/watergun/channel"
	"github.com/snaztoz/watergun/log"
	"github.com/snaztoz/watergun/user"
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
	log.Info(fmt.Sprintf("Server is listening at port %s", s.port))

	if err := s.httpServer.ListenAndServe(); !errors.Is(err, http.ErrServerClosed) {
		log.Error("Server failed to listen at the specified port", "err", err)
		return
	}
}

func (s *Server) Stop() {
	log.Info("Shutting down the server...")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := s.httpServer.Shutdown(shutdownCtx); err != nil {
		log.Error("Failed to gracefully shutting down the server", "err", err)
		return
	}

	log.Info("Server shutted down")
}

func handler() http.Handler {
	r := chi.NewRouter()

	bootstrapMiddlewares(r)
	bootstrapPublicRoutes(r)
	bootstrapAdminRoutes(r)

	return r
}

func bootstrapMiddlewares(r *chi.Mux) {
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)

	r.Use(httplog.RequestLogger(log.Logger(), &httplog.Options{
		Level:  slog.LevelInfo,
		Schema: httplog.SchemaECS,
	}))

	r.Use(middleware.Recoverer)
	r.Use(middleware.RedirectSlashes)
	r.Use(middleware.Timeout(60 * time.Second))
}

func bootstrapPublicRoutes(r *chi.Mux) {
	r.Use(middleware.Heartbeat("/up"))

	r.Get("/socket", handleWS())
}

func bootstrapAdminRoutes(r *chi.Mux) {
	r.Route("/admin", func(r chi.Router) {
		r.Use(adminRoutesAuth)

		r.Use(func(next http.Handler) http.Handler {
			return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("Content-Type", "application/json")
				next.ServeHTTP(w, r)
			})
		})

		r.Route("/channels", func(r chi.Router) {
			channelStore := channel.NewStore()
			channelDomain := channel.NewDomain(channelStore)
			channelHandler := channel.NewHandler(channelDomain)

			r.Post("/", channelHandler.CreateChannel)
			r.Get("/{id}", channelHandler.FetchChannel)

			r.Route("/{channelID}/participants", func(r chi.Router) {
				r.Get("/", channelHandler.FetchParticipantsList)
				r.Post("/", channelHandler.CreateParticipant)
			})
		})

		r.Route("/users", func(r chi.Router) {
			userStore := user.NewStore()
			userDomain := user.NewDomain(userStore)
			userHandler := user.NewHandler(userDomain)

			r.Post("/", userHandler.CreateUser)
			r.Get("/{id}", userHandler.FetchUser)
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

		if key := strings.TrimPrefix(authorization, "Bearer "); key != adminAPIKey() {
			http.Error(w, http.StatusText(403), 403)
			return
		}

		next.ServeHTTP(w, r)
	})
}
