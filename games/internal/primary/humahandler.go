package primary

import (
	"context"
	"errors"
	"github.com/danielgtaylor/huma/v2"
	games "github.com/rpatton4/mesbg-league/games/pkg"
	"github.com/rpatton4/mesbg-league/games/pkg/model"
	"github.com/rpatton4/mesbg-league/pkg/svcerrors"
	"log/slog"
)

// HumaHandler defines the HTTP handler (adapter) for Games operations received via HTTP(S).
type HumaHandler struct {
	ctrl SingleController
}

// <editor-fold desc="I/O Struct Definitions">

// Huma requires structs for both the input and output of each function registered as a handler for an HTTP operation.
// The structs most notably contain a Body field for the JSON body of the request or response, but can also contain
// headers and query parameters. Where it makes sense to use the exact models they will be used, otherwise
// local structs will be defined to hold only the fields specifically needed for the operation. In this case they will
// be in the form of a command or a query rather than a full model.
// Please keep the structs organized, with request and then response for any operation together

// GetByIDRequest defines the input for the GetByID operation.
type GetByIDRequest struct {
	// ID is the unique identifier for the game to retrieve, and it will be taken from the path with the assumption
	// that the path is set up in the form /games/{id} in the Huma API definition
	ID games.GameID `path:"id" example:"1234" doc:"The unique identifier for the game to retrieve"`
}

// GetByIDResponse defines the output for the GetByID operation.
type GetByIDResponse struct {
	// Body holds the game with the requested ID, Huma will marshall this to JSON for the HTTP response
	Body model.Game
}

// PostRequest defines the input for the Post operation, which creates a new game.
type PostRequest struct {
	// Body holds the info for the game to be created
	Body model.Game
}

type PostResponse struct {
	// Body holds the newly created game, including its assigned ID, Huma will marshall this to JSON for the HTTP response
	Body model.Game
}

type PutRequest struct {
	// ID is the unique identifier for the game to update, and it will be taken from the path with the assumption
	// that the path is set up in the form /games/{id} in the Huma API definition
	ID games.GameID `path:"id" example:"1234" doc:"The unique identifier for the game to update"`

	// Body holds the Game model to be used to update the game with the ID from the path
	Body model.Game
}

type PutResponse struct {
	// Body holds the updated game, Huma will marshall this to JSON for the HTTP response
	Body model.Game
}

type DeleteRequest struct {
	// ID is the unique identifier for the game to delete, and it will be taken from the path with the assumption
	// that the path is set up in the form /games/{id} in the Huma API definition
	ID games.GameID `path:"id" example:"1234" doc:"The unique identifier for the game to delete"`
}

//</editor-fold>

// NewHumaHandler creates a new instance of the HTTP handler for game operations.
func NewHumaHandler(c SingleController) *HumaHandler {
	return &HumaHandler{ctrl: c}
}

// GetByID queries the controller for the game with the ID taken from the path, returns it if found
// 404 is returned if no such game exists
// 400 is returned if the game ID is invalid
func (h *HumaHandler) GetByID(ctx context.Context, req *GetByIDRequest) (*GetByIDResponse, huma.StatusError) {
	slog.Debug("GetByID called", "gameID", req.ID)

	g, err := h.ctrl.GetByID(ctx, req.ID)

	if err != nil && errors.Is(err, svcerrors.ErrNotFound) {
		slog.Error("Repository does not have a game with the given ID", "gameID", req.ID, "error", err)
		return nil, huma.Error404NotFound("No such game exists")
	}

	if err != nil && errors.Is(err, svcerrors.ErrInvalidID) {
		slog.Error("The given game ID is invalid", "gameID", req.ID, "error", err)
		return nil, huma.Error400BadRequest("Invalid game ID")
	}

	return &GetByIDResponse{
		Body: *g,
	}, nil
}

// Post reads the game JSON from the HTTP call and sends it on to the controller to create the game
// 500 is returned if the game cannot be created for any reason
func (h *HumaHandler) Post(ctx context.Context, req *PostRequest) (*PostResponse, error) {
	slog.Debug("Post called", "PostRequest Body", req.Body)
	g, err := h.ctrl.Create(ctx, &req.Body)

	if err != nil {
		slog.Error("Unable to create the game", "func", "Post", "error", err)
		return nil, huma.Error500InternalServerError("Error while creating the game: " + err.Error())
	}
	slog.Debug("Created game", "game", g)
	return &PostResponse{
		Body: *g,
	}, nil
}

// Put reads the game JSON from the HTTP call and sends it on to the controller to fully update the game
// with the given ID from the path.
// 500 is returned if the game cannot be created for any reason
func (h *HumaHandler) Put(ctx context.Context, req *PutRequest) (*PutResponse, error) {
	slog.Debug("Put called", "PutRequest Body", req.Body)
	g, err := h.ctrl.Replace(ctx, &req.Body)

	if err != nil {
		slog.Error("Unable to create the game", "func", "Post", "error", err)
		return nil, huma.Error500InternalServerError("Error while creating the game: " + err.Error())
	}
	slog.Debug("Created game", "game", g)
	return &PutResponse{
		Body: *g,
	}, nil
}

// Delete deletes the game with the given ID from the path.
func (h *HumaHandler) Delete(ctx context.Context, req *DeleteRequest) (*struct{}, error) {
	slog.Debug("Delete called", "gameID", req.ID)

	_, err := h.ctrl.DeleteByID(ctx, req.ID)

	if err != nil {
		slog.Error("Controller error for game", "gameID", req.ID, "error", err)
		if errors.Is(err, svcerrors.ErrNotFound) || errors.Is(err, svcerrors.ErrInvalidID) {
			return nil, huma.Error404NotFound("No such game exists")
		}
		return nil, huma.Error500InternalServerError("Error while deleting the game: " + err.Error())
	}
	return nil, nil
}
