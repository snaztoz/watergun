package user

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"

	"github.com/snaztoz/watergun/log"
	"github.com/snaztoz/watergun/response"
)

func NewHandler(domain *domain) *handler {
	return &handler{domain: domain}
}

type handler struct {
	domain *domain
}

func (h *handler) CreateUser(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	var dto UserCreationDTO
	if err := json.NewDecoder(r.Body).Decode(&dto); err != nil {
		response.SendErrorJSON(w, err, "Failed to decode request body", 400)
		return
	}

	// TODO: validate DTO

	user, err := h.domain.createUser(dto.ID)
	if err != nil {
		response.SendErrorJSON(w, err, "Failed to create user", 422)
		return
	}

	w.WriteHeader(201)

	response.SendJSON(w, user)
}

func (h *handler) FetchUser(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	user := h.domain.fetchUser(id)
	if user == nil {
		log.Error("user does not exist", "id", id)
		http.Error(w, "User does not exist", 404)
		return
	}

	response.SendJSON(w, user)
}

type UserCreationDTO struct {
	ID string `json:"id"`
}
