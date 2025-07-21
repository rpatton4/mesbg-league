package inbound

import (
	"bytes"
	"encoding/json"
	"errors"
	"github.com/rpatton4/mesbg-league/games/pkg/header"
	"github.com/rpatton4/mesbg-league/games/pkg/model"
	"github.com/rpatton4/mesbg-league/pkg/svcerrors"
	"io"
	"log/slog"
	"net/http"
)

// Handler defines the HTTP handler for game operations.
type Handler struct {
	ctrl *Controller
}

// New creates a new instance of the HTTP handler for game operations.
func New(c *Controller) *Handler {
	return &Handler{ctrl: c}
}

// DemuxWithID takes a request with a game ID at the end of the path and routes it to the appropriate handler method.
// Assumes the path has a value called "id" which is a game ID
func (h *Handler) DemuxWithID(w http.ResponseWriter, r *http.Request) {
	id := header.GameID(r.PathValue("id"))
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

// httpGetByID queries the controller for the game with the ID taken from the path, returns it if found
// Any errors are sent directly out to the HTTP stream
func httpGetByID(h *Handler, w http.ResponseWriter, r *http.Request, id header.GameID) {
	slog.Debug("httpGetByID called", "gameID", id)

	ctx := r.Context()
	g, err := h.ctrl.GetByID(ctx, header.GameID(id))

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
		return
	}
	w.WriteHeader(http.StatusOK)
}

// httpPost reads the game json from the HTTP call and sends it on to the controller
// Any errors are sent directly out to the HTTP stream
func httpPost(h *Handler, w http.ResponseWriter, r *http.Request) {
	bodyBytes, err := io.ReadAll(r.Body)
	if err != nil {
		slog.Error("Failed to read request body", "error", err)
		http.Error(w, "Failed to read request body", http.StatusBadRequest)
		return
	}

	slog.Debug("post called with body", "body", string(bodyBytes))

	var newGame model.Game
	if err := json.NewDecoder(bytes.NewReader(bodyBytes)).Decode(&newGame); err != nil {
		slog.Error("Failed to decode newGame JSON", "error", err)
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}
	slog.Debug("Game decoded successfully", "newGame", newGame)
	createdGame, err := h.ctrl.Create(r.Context(), &newGame)

	if err != nil {
		slog.Error("Error creating newGame", "error", err)
		http.Error(w, "Error creating newGame", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
	slog.Debug("Game created successfully", "newGame", createdGame)
}

// httpPutWithID replaces the game with the given ID from the path with the one passed in.
// Any errors or ok responses are sent directly out to the HTTP stream
func httpPutWithID(h *Handler, w http.ResponseWriter, r *http.Request, id header.GameID) {
	bodyBytes, err := io.ReadAll(r.Body)
	if err != nil {
		slog.Error("Failed to read request body", "error", err)
		http.Error(w, "Failed to read request body", http.StatusBadRequest)
		return
	}

	slog.Debug("put called", "gameID", id, "body", string(bodyBytes))

	var updatedGame model.Game
	if err := json.NewDecoder(bytes.NewReader(bodyBytes)).Decode(&updatedGame); err != nil {
		slog.Error("Failed to decode updatedGame JSON", "error", err)
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}
	slog.Debug("Game decoded successfully", "updatedGame", updatedGame)

	if updatedGame.ID != "" && id != updatedGame.ID {
		slog.Error("ID in path does not match ID submitted in the game", "error", err, "pathID", id, "gameID", updatedGame.ID)
		http.Error(w, "ID in path does not match ID submitted in the game", http.StatusBadRequest)
		return
	}

	replacedGame, err := h.ctrl.Replace(r.Context(), &updatedGame)

	if err != nil {
		slog.Error("Error updating game", "error", err)
		http.Error(w, "Error updating game", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	slog.Debug("Game updated successfully", "replacedGame", replacedGame)
}

// httpDeleteByID deletes the game with the given ID from the path.
// Any errors or ok responses are sent directly out to the HTTP stream
func httpDeleteByID(h *Handler, w http.ResponseWriter, r *http.Request, id header.GameID) {
	slog.Info("httpDeleteByID called", "gameID", id)
	ok := h.ctrl.DeleteByID(r.Context(), id)
	if !ok {
		http.Error(w, "No game with that ID exists", http.StatusNotFound)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
