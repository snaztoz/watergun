package server

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"

	"github.com/snaztoz/watergun"
)

func newChannelHandler(domain *watergun.ChannelDomain) *channelHandler {
	return &channelHandler{domain: domain}
}

type channelHandler struct {
	domain *watergun.ChannelDomain
}

func (h *channelHandler) Create(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	var dto ChannelCreationDTO
	if err := json.NewDecoder(r.Body).Decode(&dto); err != nil {
		watergun.RespondWithError(w, err, "Failed to decode request body", 400)
		return
	}

	// TODO: validate DTO

	channel, err := h.domain.CreateChannel(dto.ID, dto.Name)
	if err != nil {
		watergun.RespondWithError(w, err, "Failed to create channel", 422)
		return
	}

	w.WriteHeader(201)

	if err := json.NewEncoder(w).Encode(channel); err != nil {
		watergun.RespondWithError(w, err, "Failed to write response", 500)
	}
}

func (h *channelHandler) Retrieve(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	channel := h.domain.Retrieve(id)
	if channel == nil {
		watergun.Logger().Error("Channel does not exist", "id", id)
		http.Error(w, "Channel does not exist", 404)
		return
	}

	if err := json.NewEncoder(w).Encode(channel); err != nil {
		watergun.RespondWithError(w, err, "Failed to write response", 500)
	}
}

func (h *channelHandler) CreateParticipant(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	channelID := chi.URLParam(r, "channelID")

	var dto ParticipantAdditionDTO
	if err := json.NewDecoder(r.Body).Decode(&dto); err != nil {
		watergun.RespondWithError(w, err, "Failed to decode request body", 400)
		return
	}

	// TODO: validate DTO

	participant, err := h.domain.CreateParticipant(
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

func (h *channelHandler) RetrieveParticipantsList(w http.ResponseWriter, r *http.Request) {
	channelID := chi.URLParam(r, "channelID")

	participants := h.domain.RetrieveParticipantsList(channelID)

	if err := json.NewEncoder(w).Encode(participants); err != nil {
		watergun.RespondWithError(w, err, "Failed to write response", 500)
	}
}

type ChannelCreationDTO struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type ParticipantAdditionDTO struct {
	UserID     string `json:"user_id"`
	CanPublish bool   `json:"can_publish"`
}
