package http

import (
	"bytes"
	"encoding/json"
	"errors"
	"github.com/rpatton4/mesbg-league/games/internal/controller/games"
	"github.com/rpatton4/mesbg-league/games/pkg/model"
	"github.com/rpatton4/mesbg-league/svcerrors"
	"io"
	"log/slog"
	"net/http"
)

// Handler defines the HTTP handler for game operations.
type Handler struct {
	ctrl *games.Controller
}

// New creates a new instance of the HTTP handler for game operations.
func New(c *games.Controller) *Handler {
	return &Handler{ctrl: c}
}

func (h *Handler) GetByID(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	slog.Debug("get called", "gameID", id)

	ctx := r.Context()
	g, err := h.ctrl.GetByID(ctx, model.GameID(id))

	if err != nil && errors.Is(err, svcerrors.NotFound) {
		slog.Warn("Game not found", "gameID", id)
		http.Error(w, "Game not found", http.StatusNotFound)
		return
	} else if err != nil {
		slog.Error("Repository error for game", "gameID", id, "error", err)
		http.Error(w, "Internal (repository) server error", http.StatusInternalServerError)
		return
	}

	if err := json.NewEncoder(w).Encode(g); err != nil {
		slog.Error("Failed to encode game response", "error", err)
		w.WriteHeader(http.StatusInternalServerError)
	}
}

// Demux is the HTTP handler that routes requests to the appropriate method based on the HTTP method used, with the
// exception of the GET when an ID is in the path, which is handled by the GetByID method
func (h *Handler) Demux(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		httpPost(h, w, r)
		return
	case http.MethodPut:
		httpPut(h, w, r)
		return
	case http.MethodDelete:
		httpDelete(h, w, r)
		return
	default:
		slog.Warn("Unsupported HTTP method on the path", "method", r.Method, "path", r.URL.Path)
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func httpPost(h *Handler, w http.ResponseWriter, r *http.Request) {
	bodyBytes, err := io.ReadAll(r.Body)
	if err != nil {
		slog.Error("Failed to read request body", "error", err)
		http.Error(w, "Failed to read request body", http.StatusBadRequest)
		return
	}

	slog.Debug("post called with body", "body", string(bodyBytes))

	var game model.Game
	if err := json.NewDecoder(bytes.NewReader(bodyBytes)).Decode(&game); err != nil {
		slog.Error("Failed to decode game JSON", "error", err)
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}
	slog.Debug("Game decoded successfully", "game", game)
	res, err := h.ctrl.Create(r.Context(), &game)

	if err != nil {
		slog.Error("Error creating game", "error", err)
		http.Error(w, "Error creating game", http.StatusInternalServerError)
		return
	}

	slog.Debug("Game created successfully", "game", res)
}

func httpPut(_ *Handler, _ http.ResponseWriter, _ *http.Request) {
	slog.Info("put called")
}

func httpDelete(_ *Handler, _ http.ResponseWriter, _ *http.Request) {
	slog.Info("delete called")
}
