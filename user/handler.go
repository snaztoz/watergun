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
		watergun.Logger().Error("Failed to decode request body", "err", err)
		http.Error(w, err.Error(), 400)
		return
	}

	// TODO: validate DTO

	user, err := h.domain.createUser(dto.MasterID)
	if err != nil {
		watergun.Logger().Error("Failed to create user", "err", err)
		http.Error(w, err.Error(), 422)
		return
	}

	if err := json.NewEncoder(w).Encode(user); err != nil {
		watergun.Logger().Error("Failed to write response", "err", err)
		http.Error(w, err.Error(), 500)
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
		watergun.Logger().Error("Failed to write response", "err", err)
		http.Error(w, err.Error(), 500)
	}
}

type CreationDTO struct {
	MasterID string `json:"master_id"`
}
