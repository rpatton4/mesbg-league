// Package gateway contains clients for interacting with the games service from other services.
// The package is meant to be public, and will commmit to being backwards compatible within major versions.
// The package primarily consists of the GamesGateway interface, with different implementations
// of the interface for calling it in-memory or over HTTP.
package gateway

import (
	"context"
	games "github.com/rpatton4/mesbg-league/games/pkg"
	"github.com/rpatton4/mesbg-league/games/pkg/model"
)

// GamesGateway provides a set of methods for interacting with the Games service from outside the service.
// The interface currently only provides access to single-transaction operations versus batch, but batch
// may be added in the future.
type GamesGateway interface {
	// GetByID returns the game with the given id, or a svcerrors.ErrNotFound if no game with that id exists
	GetByID(ctx context.Context, id games.GameID) (*model.Game, error)

	// Create persists a new game instance to the service and returns the game with an assigned ID.
	// A generic error is returned if the game to created is missing, while specific validation errors are
	// passed along from the service if the game is invalid in some way.
	Create(ctx context.Context, game *model.Game) (*model.Game, error)

	// Replace updates an existing game in the service with the provided game.
	// A generic error is returned if the game to replaced is not known to the service.
	// This is an idempotent operation.
	Replace(ctx context.Context, g *model.Game) (*model.Game, error)

	// DeleteByID removes the game with the given id from the service. Returns true if the game was found and
	// deleted, false otherwise. This is an idempotent operation.
	DeleteByID(ctx context.Context, id games.GameID) (bool, error)
}
