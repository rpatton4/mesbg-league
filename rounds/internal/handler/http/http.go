package http

import (
	"encoding/json"
	"errors"
	"github.com/rpatton4/mesbg-league/pkg/svcerrors"
	"github.com/rpatton4/mesbg-league/rounds/internal/domain"
	"log/slog"
	"net/http"
	"strconv"
)

// Handler defines the HTTP handler for round operations.
type Handler struct {
	ctrl *domain.Controller
}

// New creates a new instance of the HTTP handler for round operations.
func New(c *domain.Controller) *Handler {
	return &Handler{ctrl: c}
}

func (h *Handler) GetRound(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.FormValue("id"))
	if err != nil {
		slog.Error("Invalid round ID", "error", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	ctx := r.Context()
	round, err := h.ctrl.Get(ctx, id)

	if err != nil && errors.Is(err, svcerrors.ErrNotFound) {
		slog.Warn("Round not found", "roundID", id)
		w.WriteHeader(http.StatusNotFound)
		return
	} else if err != nil {
		slog.Error("Repository error for rounds", "roundID", id, "error", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if err := json.NewEncoder(w).Encode(round); err != nil {
		slog.Error("Failed to encode rounds response", "error", err)
		w.WriteHeader(http.StatusInternalServerError)
	}
}
