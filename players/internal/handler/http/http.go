package http

import (
	"encoding/json"
	"errors"
	"github.com/rpatton4/mesbg-league/players/internal/controller/players"
	"github.com/rpatton4/mesbg-league/players/pkg/model"
	"github.com/rpatton4/mesbg-league/svcerrors"
	"log/slog"
	"net/http"
)

// Handler defines the HTTP handler for players operations.
type Handler struct {
	ctrl *players.Controller
}

// New creates a new instance of the HTTP handler for players operations.
func New(c *players.Controller) *Handler {
	return &Handler{ctrl: c}
}

func (h *Handler) Demultiplex(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		getPlayer(h, w, r)
		return
	case http.MethodPost:
		postPlayer(h, w, r)
		return
	case http.MethodPut:
		putPlayer(h, w, r)
		return
	case http.MethodDelete:
		deletePlayer(h, w, r)
		return
	}
}

func getPlayer(h *Handler, w http.ResponseWriter, r *http.Request) {
	id := model.PlayerID(r.FormValue("id"))
	if id == "" {
		slog.Error("Missing players ID")
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

func postPlayer(h *Handler, w http.ResponseWriter, r *http.Request) {

}

func putPlayer(h *Handler, w http.ResponseWriter, r *http.Request) {

}

func deletePlayer(h *Handler, w http.ResponseWriter, r *http.Request) {

}
