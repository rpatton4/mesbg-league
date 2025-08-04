package http

import (
	"bytes"
	"encoding/json"
	"errors"
	"github.com/rpatton4/mesbg-league/participants/internal/controller/participants"
	"github.com/rpatton4/mesbg-league/participants/pkg/model"
	"github.com/rpatton4/mesbg-league/pkg/svcerrors"
	"io"
	"log/slog"
	"net/http"
)

// Handler defines the HTTP handler for participant operations.
type Handler struct {
	ctrl *participants.Controller
}

// New creates a new instance of the HTTP handler for participant operations.
func New(c *participants.Controller) *Handler {
	return &Handler{ctrl: c}
}

// DemuxWithID takes a request with a participant ID at the end of the path and routes it to the appropriate handler method.
// Assumes the path has a value called "id" which is a participant ID
func (h *Handler) DemuxWithID(w http.ResponseWriter, r *http.Request) {
	id := model.ParticipantID(r.PathValue("id"))
	slog.Debug("DemuxWithID called", "id", id)

	switch r.Method {
	case http.MethodGet:
		httpGetByID(h, w, r, id)
		return
	case http.MethodDelete:
		httpDeleteByID(h, w, r, id)
		return
	case http.MethodPut:
		httpPutWithID(h, w, r, id)
		return
	default:
		slog.Warn("Unsupported HTTP method on the path", "method", r.Method, "path", r.URL.Path)
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

// Demux is the HTTP handler that routes requests to the appropriate method based on the HTTP method used, with the
// exception of calls when an ID is in the path, which are handled by the DemuxWithID method
func (h *Handler) Demux(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		httpPost(h, w, r)
		return
	default:
		slog.Warn("Unsupported HTTP method on the path", "method", r.Method, "path", r.URL.Path)
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func httpGetByID(h *Handler, w http.ResponseWriter, r *http.Request, id model.ParticipantID) {
	slog.Debug("httpGetByID called", "participantID", id)

	ctx := r.Context()
	g, err := h.ctrl.GetByID(ctx, model.ParticipantID(id))

	if err != nil && errors.Is(err, svcerrors.ErrNotFound) {
		slog.Warn("Participant not found", "participantID", id)
		http.Error(w, "Participant not found", http.StatusNotFound)
		return
	} else if err != nil {
		slog.Error("Repository error for participant", "participantID", id, "error", err)
		http.Error(w, "Internal (repository) server error", http.StatusInternalServerError)
		return
	}

	if err := json.NewEncoder(w).Encode(g); err != nil {
		slog.Error("Failed to encode participant response", "error", err)
		w.WriteHeader(http.StatusInternalServerError)
	}
	w.WriteHeader(http.StatusOK)
}

func httpPost(h *Handler, w http.ResponseWriter, r *http.Request) {
	bodyBytes, err := io.ReadAll(r.Body)
	if err != nil {
		slog.Error("Failed to read request body", "error", err)
		http.Error(w, "Failed to read request body", http.StatusBadRequest)
		return
	}

	slog.Debug("post called with body", "body", string(bodyBytes))

	var newParticipant model.Participant
	if err := json.NewDecoder(bytes.NewReader(bodyBytes)).Decode(&newParticipant); err != nil {
		slog.Error("Failed to decode participant JSON", "error", err)
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}
	slog.Debug("Participant decoded successfully", "game", newParticipant)
	createdParticipant, err := h.ctrl.Create(r.Context(), &newParticipant)

	if err != nil {
		slog.Error("Error creating participant", "error", err)
		http.Error(w, "Error creating participant", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	slog.Debug("Participant created successfully", "participant", createdParticipant)
}

// httpPutWithID replaces the participant with the given ID from the path with the one passed in.
// Any errors or ok responses are sent directly out to the HTTP stream
func httpPutWithID(h *Handler, w http.ResponseWriter, r *http.Request, id model.ParticipantID) {
	bodyBytes, err := io.ReadAll(r.Body)
	if err != nil {
		slog.Error("Failed to read request body", "error", err)
		http.Error(w, "Failed to read request body", http.StatusBadRequest)
		return
	}

	slog.Debug("put called", "participantID", id, "body", string(bodyBytes))

	var updatedParticipant model.Participant
	if err := json.NewDecoder(bytes.NewReader(bodyBytes)).Decode(&updatedParticipant); err != nil {
		slog.Error("Failed to decode updatedParticipant JSON", "error", err)
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}
	slog.Debug("Participant decoded successfully", "updatedParticipant", updatedParticipant)

	if updatedParticipant.ID != "" && id != updatedParticipant.ID {
		slog.Error("ID in path does not match ID submitted in the participant", "error", err, "pathID", id, "participantID", updatedParticipant.ID)
		http.Error(w, "ID in path does not match ID submitted in the participant", http.StatusBadRequest)
		return
	}

	replacedParticipant, err := h.ctrl.Replace(r.Context(), &updatedParticipant)

	if err != nil {
		slog.Error("Error updating participant", "error", err)
		http.Error(w, "Error updating participant", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	slog.Debug("Participant updated successfully", "replacedPlayer", replacedParticipant)
}

// httpDeleteByID deletes the participant with the given ID from the path.
// Any errors or ok responses are sent directly out to the HTTP stream
func httpDeleteByID(h *Handler, w http.ResponseWriter, r *http.Request, id model.ParticipantID) {
	slog.Info("httpDeleteByID called", "participantID", id)
	ok := h.ctrl.DeleteByID(r.Context(), id)
	if !ok {
		http.Error(w, "No participant with that ID exists", http.StatusNotFound)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
