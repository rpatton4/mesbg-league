package primary

import (
	"context"
	"errors"
	"github.com/rpatton4/mesbg-league/games/internal/secondary"
	"github.com/rpatton4/mesbg-league/games/pkg"
	"github.com/rpatton4/mesbg-league/games/pkg/model"
	"github.com/rpatton4/mesbg-league/pkg/svcerrors"
)

// TxnController implements the single controller for game operations.
type TxnController struct {
	repo secondary.Repository
}

// NewTxnController creates a new instance of the games controller for transactional behavior in the sense of realtime
// operations on a game, versus batch
func NewTxnController(r secondary.Repository) *TxnController {
	return &TxnController{repo: r}
}

// GetByID returns the game with the given id, or a svcerrors.ErrNotFound if no game with that id exists
func (c *TxnController) GetByID(ctx context.Context, id pkg.GameID) (*model.Game, error) {
	return c.repo.GetByID(ctx, id)
}

// Create persists a new game instance to the repository and returns the game with an assigned ID.
// A generic error is returned if the game to created is missing, while specific validation errors are
// passed along from the repository if the game is invalid in some way.
func (c *TxnController) Create(ctx context.Context, g *model.Game) (*model.Game, error) {
	if g == nil {
		return nil, errors.New("the game to be created cannot be nil")
	}
	return c.repo.Create(ctx, g)
}

// Replace updates an existing game in the repository with the provided game.
// A generic error is returned if the game to replaced is not present in the data store.
func (c *TxnController) Replace(ctx context.Context, g *model.Game) (*model.Game, error) {
	if g == nil {
		return nil, errors.New("the game to be replaced cannot be nil")
	}
	return c.repo.Replace(ctx, g)
}

// DeleteByID removes the game with the given id from the repository. Returns true if the game was found and
// deleted, false otherwise. This is an idempotent operation.
func (c *TxnController) DeleteByID(ctx context.Context, id pkg.GameID) (bool, error) {
	if id == "" {
		return false, svcerrors.ErrInvalidID
	}
	return c.repo.DeleteByID(ctx, id)
}
