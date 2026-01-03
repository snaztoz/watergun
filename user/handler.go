package user

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"

	"github.com/snaztoz/watergun"
)

func NewHandler(domain *domain) *handler {
	return &handler{domain: domain}
}

type handler struct {
	domain *domain
}

func (h *handler) CreateUser(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	var dto CreationDTO
	if err := json.NewDecoder(r.Body).Decode(&dto); err != nil {
		watergun.RespondWithError(w, err, "Failed to decode request body", 400)
		return
	}

	// TODO: validate DTO

	user, err := h.domain.createUser(dto.ID)
	if err != nil {
		watergun.RespondWithError(w, err, "Failed to create user", 422)
		return
	}

	w.WriteHeader(201)

	if err := json.NewEncoder(w).Encode(user); err != nil {
		watergun.RespondWithError(w, err, "Failed to write response", 500)
	}
}

func (h *handler) RetrieveUser(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	user := h.domain.retrieveUser(id)
	if user == nil {
		watergun.Logger().Error("User does not exist", "id", id)
		http.Error(w, "User does not exist", 404)
		return
	}

	if err := json.NewEncoder(w).Encode(user); err != nil {
		watergun.RespondWithError(w, err, "Failed to write response", 500)
	}
}

type CreationDTO struct {
	ID string `json:"id"`
}
