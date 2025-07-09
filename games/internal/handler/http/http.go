package http

import (
	"encoding/json"
	"errors"
	"github.com/rpatton4/mesbg-league/games/internal/controller/games"
	"github.com/rpatton4/mesbg-league/svcerrors"
	"log/slog"
	"net/http"
	"strconv"
)

// Handler defines the HTTP handler for game operations.
type Handler struct {
	ctrl *games.Controller
}

// New creates a new instance of the HTTP handler for game operations.
func New(c *games.Controller) *Handler {
	return &Handler{ctrl: c}
}

func (h *Handler) GetGame(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.FormValue("id"))
	if err != nil {
		slog.Error("Invalid game ID", "error", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	ctx := r.Context()
	p, err := h.ctrl.Get(ctx, id)

	if err != nil && errors.Is(err, svcerrors.NotFound) {
		slog.Warn("Game not found", "gameID", id)
		w.WriteHeader(http.StatusNotFound)
		return
	} else if err != nil {
		slog.Error("Repository error for game", "gameID", id, "error", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if err := json.NewEncoder(w).Encode(p); err != nil {
		slog.Error("Failed to encode game response", "error", err)
		w.WriteHeader(http.StatusInternalServerError)
	}
}
