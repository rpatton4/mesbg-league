package http

import (
	"encoding/json"
	"errors"
	"github.com/rpatton4/mesbg-league/players/internal/controller/player"
	"github.com/rpatton4/mesbg-league/svcerrors"
	"log/slog"
	"net/http"
	"strconv"
)

// Handler defines the HTTP handler for players operations.
type Handler struct {
	ctrl *player.Controller
}

// New creates a new instance of the HTTP handler for players operations.
func New(c *player.Controller) *Handler {
	return &Handler{ctrl: c}
}

func (h *Handler) GetPlayer(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.FormValue("id"))
	if err != nil {
		slog.Error("Invalid players ID", "error", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	ctx := r.Context()
	p, err := h.ctrl.Get(ctx, id)

	if err != nil && errors.Is(err, svcerrors.NotFound) {
		slog.Warn("Player not found", "playerID", id)
		w.WriteHeader(http.StatusNotFound)
		return
	} else if err != nil {
		slog.Error("Repository error for players", "playerID", id, "error", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if err := json.NewEncoder(w).Encode(p); err != nil {
		slog.Error("Failed to encode players response", "error", err)
		w.WriteHeader(http.StatusInternalServerError)
	}
}
