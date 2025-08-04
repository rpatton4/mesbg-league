package http

import (
	"bytes"
	"encoding/json"
	"errors"
	"github.com/rpatton4/mesbg-league/pkg/svcerrors"
	ctrl "github.com/rpatton4/mesbg-league/players/internal/controller/players"
	players "github.com/rpatton4/mesbg-league/players/pkg"
	"github.com/rpatton4/mesbg-league/players/pkg/model"
	"io"
	"log/slog"
	"net/http"
)

// Handler defines the HTTP handler for players operations.
type Handler struct {
	ctrl *ctrl.Controller
}

// New creates a new instance of the HTTP handler for players operations.
func New(c *ctrl.Controller) *Handler {
	return &Handler{ctrl: c}
}

// DemuxWithID takes a request with a player ID at the end of the path and routes it to the appropriate handler method.
// Assumes the path has a value called "id" which is a player ID
func (h *Handler) DemuxWithID(w http.ResponseWriter, r *http.Request) {
	id := players.PlayerID(r.PathValue("id"))
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

func httpGetByID(h *Handler, w http.ResponseWriter, r *http.Request, id players.PlayerID) {
	slog.Debug("httpGetByID called", "playerID", id)

	ctx := r.Context()
	g, err := h.ctrl.GetByID(ctx, players.PlayerID(id))

	if err != nil && errors.Is(err, svcerrors.ErrNotFound) {
		slog.Warn("Player not found", "playerID", id)
		http.Error(w, "Player not found", http.StatusNotFound)
		return
	} else if err != nil {
		slog.Error("Repository error for player", "playerID", id, "error", err)
		http.Error(w, "Internal (repository) server error", http.StatusInternalServerError)
		return
	}

	if err := json.NewEncoder(w).Encode(g); err != nil {
		slog.Error("Failed to encode player response", "error", err)
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

	var newPlayer model.Player
	if err := json.NewDecoder(bytes.NewReader(bodyBytes)).Decode(&newPlayer); err != nil {
		slog.Error("Failed to decode player JSON", "error", err)
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}
	slog.Debug("Player decoded successfully", "game", newPlayer)
	createdPlayer, err := h.ctrl.Create(r.Context(), &newPlayer)

	if err != nil {
		slog.Error("Error creating player", "error", err)
		http.Error(w, "Error creating player", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	slog.Debug("Player created successfully", "player", createdPlayer)
}

// httpPutWithID replaces the player with the given ID from the path with the one passed in.
// Any errors or ok responses are sent directly out to the HTTP stream
func httpPutWithID(h *Handler, w http.ResponseWriter, r *http.Request, id players.PlayerID) {
	bodyBytes, err := io.ReadAll(r.Body)
	if err != nil {
		slog.Error("Failed to read request body", "error", err)
		http.Error(w, "Failed to read request body", http.StatusBadRequest)
		return
	}

	slog.Debug("put called", "playerID", id, "body", string(bodyBytes))

	var updatedPlayer model.Player
	if err := json.NewDecoder(bytes.NewReader(bodyBytes)).Decode(&updatedPlayer); err != nil {
		slog.Error("Failed to decode updatedPlayer JSON", "error", err)
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}
	slog.Debug("Player decoded successfully", "updatedPlayer", updatedPlayer)

	if updatedPlayer.ID != "" && id != updatedPlayer.ID {
		slog.Error("ID in path does not match ID submitted in the player", "error", err, "pathID", id, "gameID", updatedPlayer.ID)
		http.Error(w, "ID in path does not match ID submitted in the player", http.StatusBadRequest)
		return
	}

	replacedPlayer, err := h.ctrl.Replace(r.Context(), &updatedPlayer)

	if err != nil {
		slog.Error("Error updating player", "error", err)
		http.Error(w, "Error updating player", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	slog.Debug("Player updated successfully", "replacedPlayer", replacedPlayer)
}

// httpDeleteByID deletes the player with the given ID from the path.
// Any errors or ok responses are sent directly out to the HTTP stream
func httpDeleteByID(h *Handler, w http.ResponseWriter, r *http.Request, id players.PlayerID) {
	slog.Info("httpDeleteByID called", "playerID", id)
	ok := h.ctrl.DeleteByID(r.Context(), id)
	if !ok {
		http.Error(w, "No player with that ID exists", http.StatusNotFound)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
