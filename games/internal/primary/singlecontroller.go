package primary

import (
	"context"
	"github.com/rpatton4/mesbg-league/games/internal/secondary"
	games "github.com/rpatton4/mesbg-league/games/pkg"
	"github.com/rpatton4/mesbg-league/games/pkg/model"
)

//go:generate mockgen --destination ./mocks/controller.go github.com/rpatton4/mesbg-league/games/internal/primary SingleController
type SingleController interface {
	// GetByID returns the game with the given id, or a svcerrors.ErrNotFound if no game with that id exists
	GetByID(ctx context.Context, id games.GameID) (*model.Game, error)

	// Create persists a new game instance to the repository and returns the game with an assigned ID.
	// A generic error is returned if the game to created is missing, while specific validation errors are
	// passed along from the repository if the game is invalid in some way.
	Create(ctx context.Context, g *model.Game) (*model.Game, error)

	// Replace updates an existing game in the repository with the provided game.
	// A generic error is returned if the game to replaced is not present in the data store.
	Replace(ctx context.Context, g *model.Game) (*model.Game, error)

	// DeleteByID removes the game with the given id from the repository. Returns true if the game was found and
	// deleted, false otherwise. This is an idempotent operation.
	DeleteByID(ctx context.Context, id games.GameID) (bool, error)
}

// NewDefaultSingleController creates an instance of the default single controller implementation. This default is controlled
// by configuration.
// TODO: Implement the configuration control for this
func NewDefaultSingleController(repo secondary.Repository) SingleController {
	return &TxnController{
		repo: repo,
	}
}
