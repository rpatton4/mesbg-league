package gamesprimaryadapters

import (
	"bytes"
	"context"
	"encoding/json"
	"github.com/danielgtaylor/huma/v2"
	"github.com/rpatton4/mesbg-league/games/internal/gamesprimaryports"
	"github.com/rpatton4/mesbg-league/games/pkg/header"
	"github.com/rpatton4/mesbg-league/games/pkg/model"
	"io"
	"log/slog"
	"net/http"
)

// HumaHandler defines the HTTP handler (adapter) for Games operations received via HTTP(S).
type HumaHandler struct {
	ctrl *gamesprimaryports.TxnController
}

// Huma requires structs for both the input and output of each function registered as a handler for an HTTP operation.
// The structs

type GetByIDRequest struct {
	// ID is the unique identifier for the game to retrieve.
	ID header.GameID `path:"id" example:"1234" doc:"The unique identifier for the game to retrieve"`
}

type GetByIDResponse struct {
	Body model.Game
}

type PostRequest struct {
	Body GameCreateCommand
}

type GameCreateCommand struct {
	// Side1ID is the identifier of the player for the first side in the game.
	// First does not imply that this player acts first, it is simply a designator
	Side1ID string `json:"side1Id" example:"5678" doc:"The unique identifier for the first side in the game"`

	// Side2ID is the identifier of the player for the second side in the game.
	// Second does not imply that this player acts second, it is simply a designator
	Side2ID string `json:"side2Id" example:"9999" doc:"The unique identifier for the second side in the game"`

	// RoundID is the key to the round in a league which this game is part of
	RoundID string `json:"roundId" example:"9876" doc:"The unique identifier for the round associated with this game, if there is one"`

	// Side1TotalVictoryPoints is the total victory points scored by the first side in the game
	Side1TotalVictoryPoints int `json:"side1TotalVictoryPoints,omitempty" example:"12" doc:"The total number of victory points scored by the first player"`

	// Side2TotalVictoryPoints is the total victory points scored by the second side in the game
	Side2TotalVictoryPoints int `json:"side2TotalVictoryPoints,omitempty" example:"4" doc:"The total number of victory points second by the first player"`

	// Side1TotalGeneralsKilled is true if the side 1 player killed the opposing general
	Side1KilledGeneral bool `json:"side1KilledGeneral,omitempty" example:"true" doc:"True if the first player killed the opposing general, false otherwise"`

	// Side2TotalGeneralsKilled is true if the side 2 player killed the opposing general
	Side2KilledGeneral bool `json:"side2KilledGeneral,omitempty" example:"false" doc:"True if the second player killed the opposing general, false otherwise"`

	// Status is used to track whether the game is scheduled, played, conceded etc.
	// See the GameStateXYZ constants for potential values.
	Status model.GameState `json:"status,omitempty" example:"1" doc:"The current state of the game, indicating whether it is scheduled, in progress, completed etc."`
}

type PostResponse struct {
	Body model.Game
}

// NewHTTPHandler creates a new instance of the HTTP handler for game operations.
func NewHTTPHandler(c *gamesprimaryports.TxnController) *HumaHandler {
	return &HumaHandler{ctrl: c}
}

// HumaGetByID queries the controller for the game with the ID taken from the path, returns it if found
// Any errors are sent directly out to the HTTP stream
func (h *HumaHandler) HumaGetByID(ctx context.Context, req *GetByIDRequest) (*GetByIDResponse, error) {
	slog.Debug("humaGetByID called", "gameID", req.ID)

	g, err := h.ctrl.GetByID(ctx, req.ID)

	if err != nil {
		slog.Error("Repository error for game", "gameID", req.ID, "error", err)
		return nil, huma.Error404NotFound("No such game exists")
	}

	return &GetByIDResponse{
		Body: *g,
	}, nil
}

// HumaPost reads the game json from the HTTP call and sends it on to the controller
// Any errors are sent directly out to the HTTP stream
func (h *HumaHandler) HumaPost(ctx context.Context, req *PostRequest) (*PostResponse, error) {
	slog.Debug("HumaPost called", "PostRequest Body", req.Body)
	g, err := h.ctrl.Create(ctx, &model.Game{
		Side1ID:                 req.Body.Side1ID,
		Side2ID:                 req.Body.Side2ID,
		RoundID:                 req.Body.RoundID,
		Side1TotalVictoryPoints: req.Body.Side1TotalVictoryPoints,
		Side2TotalVictoryPoints: req.Body.Side2TotalVictoryPoints,
		Side1KilledGeneral:      req.Body.Side1KilledGeneral,
		Side2KilledGeneral:      req.Body.Side2KilledGeneral,
		Status:                  req.Body.Status,
	})

	if err != nil {
		slog.Error("Unable to create the game", "func", "HumaPost", "error", err)
		return nil, huma.Error500InternalServerError("Error while creating the game: " + err.Error())
	}
	slog.Debug("Created game", "game", g)
	return &PostResponse{
		Body: *g,
	}, nil

}

// httpPutWithID replaces the game with the given ID from the path with the one passed in.
// Any errors or ok responses are sent directly out to the HTTP stream
func httpPutWithID(h *HumaHandler, w http.ResponseWriter, r *http.Request, id header.GameID) {
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
func httpDeleteByID(h *HumaHandler, w http.ResponseWriter, r *http.Request, id header.GameID) {
	slog.Info("httpDeleteByID called", "gameID", id)
	ok := h.ctrl.DeleteByID(r.Context(), id)
	if !ok {
		http.Error(w, "No game with that ID exists", http.StatusNotFound)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
