package channel

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

func (h *handler) CreateChannel(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	var dto ChannelCreationDTO
	if err := json.NewDecoder(r.Body).Decode(&dto); err != nil {
		response.SendErrorJSON(w, err, "Failed to decode request body", 400)
		return
	}

	// TODO: validate DTO

	channel, err := h.domain.createChannel(dto.ID, dto.Name)
	if err != nil {
		response.SendErrorJSON(w, err, "Failed to create channel", 422)
		return
	}

	w.WriteHeader(201)

	response.SendJSON(w, channel)
}

func (h *handler) FetchChannel(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	channel := h.domain.fetchChannel(id)
	if channel == nil {
		log.Error("Channel does not exist", "id", id)
		http.Error(w, "Channel does not exist", 404)
		return
	}

	response.SendJSON(w, channel)
}

func (h *handler) CreateParticipant(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	channelID := chi.URLParam(r, "channelID")

	var dto ParticipantCreationDTO
	if err := json.NewDecoder(r.Body).Decode(&dto); err != nil {
		response.SendErrorJSON(w, err, "Failed to decode request body", 400)
		return
	}

	// TODO: validate DTO

	participant, err := h.domain.createParticipant(
		channelID,
		dto.UserID,
		dto.CanPublish,
	)
	if err != nil {
		response.SendErrorJSON(w, err, "Failed to create channel participant", 422)
		return
	}

	w.WriteHeader(201)

	response.SendJSON(w, participant)
}

func (h *handler) FetchParticipantsList(w http.ResponseWriter, r *http.Request) {
	channelID := chi.URLParam(r, "channelID")

	participants := h.domain.fetchParticipantsList(channelID)

	response.SendJSON(w, participants)
}

type ChannelCreationDTO struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type ParticipantCreationDTO struct {
	UserID     string `json:"user_id"`
	CanPublish bool   `json:"can_publish"`
}
