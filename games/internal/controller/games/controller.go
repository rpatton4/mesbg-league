package games

import (
	"context"
	"github.com/rpatton4/mesbg-league/games/pkg/model"
)

type gamesRepository interface {
	Get(ctx context.Context, id int) (*model.Game, error)
}

// Controller defines the simple controller for game operations.
type Controller struct {
	repo gamesRepository
}

// New creates a new instance of the games controller.
func New(r gamesRepository) *Controller {
	return &Controller{repo: r}
}

// Get returns the game with the given id, or a svcerrors.NotFound if no game with that id exists
func (c *Controller) Get(ctx context.Context, id int) (*model.Game, error) {
	return c.repo.Get(ctx, id)
}
