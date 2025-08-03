package secondary

import (
	"context"
	"github.com/rpatton4/mesbg-league/games/pkg"
	"github.com/rpatton4/mesbg-league/games/pkg/model"
)

// Repository defines the port for writing Games to persistent storage
type Repository interface {
	// GetByID retrieves a game by ID from the repository, if no game with the given
	// ID exists, it returns nil, svcerrors.ErrNotFound.
	GetByID(ctx context.Context, id pkg.GameID) (*model.Game, error)

	// Create persists a new game instance to the repository and returns the game with an assigned ID.
	Create(ctx context.Context, g *model.Game) (*model.Game, error)

	// Replace completely replaces an existing game instance with the provided one, using the ID from the provided game
	// to find which game to replace. This cannot be used to create a new Game, and it is an idempotent operation.
	// This is an intended equivalent to the HTTP PUT operation, though it purposefully does not allow the create which
	// PUT is sometimes interpreted as allowing (because that leaves ID creation up to the client).
	Replace(ctx context.Context, g *model.Game) (*model.Game, error)

	// DeleteByID deletes an existing game instance in the repository. Returns true if the game was found and
	// deleted, false otherwise. This is an idempotent operation.
	DeleteByID(ctx context.Context, id pkg.GameID) (bool, error)
}

// NewDefaultRepository creates an instance of the default repository implementation. This default is controlled
// by configuration.
// TODO: Implement the configuration control for this
func NewDefaultRepository() Repository {
	return NewMemoryRepository()
}
