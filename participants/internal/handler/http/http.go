package http

import (
	"encoding/json"
	"errors"
	"github.com/rpatton4/mesbg-league/participants/internal/controller/participants"
	"github.com/rpatton4/mesbg-league/svcerrors"
	"log/slog"
	"net/http"
	"strconv"
)

// Handler defines the HTTP handler for participant operations.
type Handler struct {
	ctrl *participants.Controller
}

// New creates a new instance of the HTTP handler for participant operations.
func New(c *participants.Controller) *Handler {
	return &Handler{ctrl: c}
}

func (h *Handler) GetParticipant(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.FormValue("id"))
	if err != nil {
		slog.Error("Invalid participant ID", "error", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	ctx := r.Context()
	p, err := h.ctrl.Get(ctx, id)

	if err != nil && errors.Is(err, svcerrors.NotFound) {
		slog.Warn("Participant not found", "participantID", id)
		w.WriteHeader(http.StatusNotFound)
		return
	} else if err != nil {
		slog.Error("Repository error for participant", "participantID", id, "error", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if err := json.NewEncoder(w).Encode(p); err != nil {
		slog.Error("Failed to encode participant response", "error", err)
		w.WriteHeader(http.StatusInternalServerError)
	}
}
