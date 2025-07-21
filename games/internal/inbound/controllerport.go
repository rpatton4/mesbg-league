package inbound

import (
	"context"
	"errors"
	"github.com/rpatton4/mesbg-league/games/pkg/header"
	"github.com/rpatton4/mesbg-league/games/pkg/model"
)

type gamesRepository interface {
	GetByID(ctx context.Context, id header.GameID) (*model.Game, error)
	Create(ctx context.Context, g *model.Game) (*model.Game, error)
	Replace(ctx context.Context, g *model.Game) (*model.Game, error)
	DeleteByID(ctx context.Context, id header.GameID) bool
}

// Controller defines the simple controller for game operations.
type Controller struct {
	repo gamesRepository
}

// New creates a new instance of the games controller.
func New(r gamesRepository) *Controller {
	return &Controller{repo: r}
}

// GetByID returns the game with the given id, or a svcerrors.NotFound if no game with that id exists
func (c *Controller) GetByID(ctx context.Context, id header.GameID) (*model.Game, error) {
	return c.repo.GetByID(ctx, id)
}

// Create persists a new game instance to the repository and returns the game with an assigned ID.
// A generic error is returned if the game to created is missing, while specific validation errors are
// passed along from the repository if the game is invalid in some way.
func (c *Controller) Create(ctx context.Context, g *model.Game) (*model.Game, error) {
	if g == nil {
		return nil, errors.New("the game to be created cannot be nil")
	}
	return c.repo.Create(ctx, g)
}

// Replace updates an existing game in the repository with the provided game.
// A generic error is returned if the game to replaced is not present in the data store.
func (c *Controller) Replace(ctx context.Context, g *model.Game) (*model.Game, error) {
	if g == nil {
		return nil, errors.New("the game to be replaced cannot be nil")
	}
	return c.repo.Replace(ctx, g)
}

// DeleteByID removes the game with the given id from the repository. Returns true if the game was found and
// deleted, false otherwise. This is an idempotent operation.
func (c *Controller) DeleteByID(ctx context.Context, id header.GameID) bool {
	if id == "" {
		return false
	}
	return c.repo.DeleteByID(ctx, id)
}
