package channel

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

func (h *handler) CreateChannel(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	var dto CreationDTO
	if err := json.NewDecoder(r.Body).Decode(&dto); err != nil {
		watergun.Logger().Error("Failed to decode request body", "err", err)
		http.Error(w, err.Error(), 400)
		return
	}

	// TODO: validate DTO

	channel, err := h.domain.createChannel(dto.ID, dto.Name)
	if err != nil {
		watergun.Logger().Error("Failed to create channel", "err", err)
		http.Error(w, err.Error(), 422)
		return
	}

	if err := json.NewEncoder(w).Encode(channel); err != nil {
		watergun.Logger().Error("Failed to write response", "err", err)
		http.Error(w, err.Error(), 500)
	}
}

func (h *handler) RetrieveChannel(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	channel := h.domain.retrieveChannel(id)
	if channel == nil {
		watergun.Logger().Error("Channel does not exist", "id", id)
		http.Error(w, "Channel does not exist", 404)
		return
	}

	if err := json.NewEncoder(w).Encode(channel); err != nil {
		watergun.Logger().Error("Failed to write response", "err", err)
		http.Error(w, err.Error(), 500)
	}
}

type CreationDTO struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}
