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

func (h *channelHandler) createChannel(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	var dto ChannelCreationDTO
	if err := json.NewDecoder(r.Body).Decode(&dto); err != nil {
		respondWithError(w, err, "Failed to decode request body", 400)
		return
	}

	// TODO: validate DTO

	channel, err := h.domain.CreateChannel(dto.ID, dto.Name)
	if err != nil {
		respondWithError(w, err, "Failed to create channel", 422)
		return
	}

	w.WriteHeader(201)

	respondWithJSON(w, channel)
}

func (h *channelHandler) fetchChannel(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	channel := h.domain.Fetch(id)
	if channel == nil {
		watergun.Logger().Error("Channel does not exist", "id", id)
		http.Error(w, "Channel does not exist", 404)
		return
	}

	respondWithJSON(w, channel)
}

func (h *channelHandler) createParticipant(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	channelID := chi.URLParam(r, "channelID")

	var dto ParticipantCreationDTO
	if err := json.NewDecoder(r.Body).Decode(&dto); err != nil {
		respondWithError(w, err, "Failed to decode request body", 400)
		return
	}

	// TODO: validate DTO

	participant, err := h.domain.CreateParticipant(
		channelID,
		dto.UserID,
		dto.CanPublish,
	)
	if err != nil {
		respondWithError(w, err, "Failed to create channel participant", 422)
		return
	}

	w.WriteHeader(201)

	respondWithJSON(w, participant)
}

func (h *channelHandler) fetchParticipantsList(w http.ResponseWriter, r *http.Request) {
	channelID := chi.URLParam(r, "channelID")

	participants := h.domain.FetchParticipantsList(channelID)

	respondWithJSON(w, participants)
}

type ChannelCreationDTO struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type ParticipantCreationDTO struct {
	UserID     string `json:"user_id"`
	CanPublish bool   `json:"can_publish"`
}
