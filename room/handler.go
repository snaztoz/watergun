package room

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

func (h *handler) CreateRoom(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	var dto RoomCreationDTO
	if err := json.NewDecoder(r.Body).Decode(&dto); err != nil {
		response.SendErrorJSON(w, err, "Failed to decode request body", 400)
		return
	}

	// TODO: validate DTO

	room, err := h.domain.createRoom(dto.ID, dto.Name)
	if err != nil {
		response.SendErrorJSON(w, err, "Failed to create room", 422)
		return
	}

	w.WriteHeader(201)

	response.SendJSON(w, room)
}

func (h *handler) FetchRoom(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	room := h.domain.fetchRoom(id)
	if room == nil {
		log.Error("Room does not exist", "id", id)
		http.Error(w, "Room does not exist", 404)
		return
	}

	response.SendJSON(w, room)
}

func (h *handler) CreateParticipant(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	roomID := chi.URLParam(r, "roomID")

	var dto ParticipantCreationDTO
	if err := json.NewDecoder(r.Body).Decode(&dto); err != nil {
		response.SendErrorJSON(w, err, "Failed to decode request body", 400)
		return
	}

	// TODO: validate DTO

	participant, err := h.domain.createParticipant(
		roomID,
		dto.UserID,
		dto.CanPublish,
	)
	if err != nil {
		response.SendErrorJSON(w, err, "Failed to create room participant", 422)
		return
	}

	w.WriteHeader(201)

	response.SendJSON(w, participant)
}

func (h *handler) FetchParticipantsList(w http.ResponseWriter, r *http.Request) {
	roomID := chi.URLParam(r, "roomID")

	participants := h.domain.fetchParticipantsList(roomID)

	response.SendJSON(w, participants)
}

type RoomCreationDTO struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type ParticipantCreationDTO struct {
	UserID     string `json:"user_id"`
	CanPublish bool   `json:"can_publish"`
}
