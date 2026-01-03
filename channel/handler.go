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

func (h *handler) Create(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	var dto CreationDTO
	if err := json.NewDecoder(r.Body).Decode(&dto); err != nil {
		watergun.RespondWithError(w, err, "Failed to decode request body", 400)
		return
	}

	// TODO: validate DTO

	channel, err := h.domain.create(dto.ID, dto.Name)
	if err != nil {
		watergun.RespondWithError(w, err, "Failed to create channel", 422)
		return
	}

	w.WriteHeader(201)

	if err := json.NewEncoder(w).Encode(channel); err != nil {
		watergun.RespondWithError(w, err, "Failed to write response", 500)
	}
}

func (h *handler) Retrieve(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	channel := h.domain.retrieve(id)
	if channel == nil {
		watergun.Logger().Error("Channel does not exist", "id", id)
		http.Error(w, "Channel does not exist", 404)
		return
	}

	if err := json.NewEncoder(w).Encode(channel); err != nil {
		watergun.RespondWithError(w, err, "Failed to write response", 500)
	}
}

func (h *handler) CreateParticipant(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	channelID := chi.URLParam(r, "channelID")

	var dto ParticipantAdditionDTO
	if err := json.NewDecoder(r.Body).Decode(&dto); err != nil {
		watergun.RespondWithError(w, err, "Failed to decode request body", 400)
		return
	}

	// TODO: validate DTO

	participant, err := h.domain.createParticipant(
		channelID,
		dto.UserID,
		dto.CanPublish,
	)
	if err != nil {
		watergun.RespondWithError(w, err, "Failed to create channel participant", 422)
		return
	}

	w.WriteHeader(201)

	if err := json.NewEncoder(w).Encode(participant); err != nil {
		watergun.RespondWithError(w, err, "Failed to write response", 500)
		return
	}
}

func (h *handler) RetrieveParticipantsList(w http.ResponseWriter, r *http.Request) {
	channelID := chi.URLParam(r, "channelID")

	participants := h.domain.retrieveParticipantsList(channelID)

	if err := json.NewEncoder(w).Encode(participants); err != nil {
		watergun.RespondWithError(w, err, "Failed to write response", 500)
	}
}

type CreationDTO struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type ParticipantAdditionDTO struct {
	UserID     string `json:"user_id"`
	CanPublish bool   `json:"can_publish"`
}
