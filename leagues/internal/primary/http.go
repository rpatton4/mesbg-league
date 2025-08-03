package primary

import (
	"encoding/json"
	"errors"
	"github.com/rpatton4/mesbg-league/pkg/svcerrors"
	"log/slog"
	"net/http"
	"strconv"
)

// Handler defines the HTTP handler for league operations.
type Handler struct {
	ctrl *Controller
}

// New creates a new instance of the HTTP handler for league operations.
func New(c *Controller) *Handler {
	return &Handler{ctrl: c}
}

func (h *Handler) GetLeague(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.FormValue("id"))
	if err != nil {
		slog.Error("Invalid league ID", "error", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	ctx := r.Context()
	p, err := h.ctrl.Get(ctx, id)

	if err != nil && errors.Is(err, svcerrors.ErrNotFound) {
		slog.Warn("League not found", "leagueID", id)
		w.WriteHeader(http.StatusNotFound)
		return
	} else if err != nil {
		slog.Error("Repository error for league", "leagueID", id, "error", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if err := json.NewEncoder(w).Encode(p); err != nil {
		slog.Error("Failed to encode league response", "error", err)
		w.WriteHeader(http.StatusInternalServerError)
	}
}
