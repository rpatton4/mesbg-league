package player

import (
	"context"
	"github.com/rpatton4/mesbg-league/player/pkg/model"
)

type playerRepository interface {
	Get(ctx context.Context, id int) (*model.Player, error)
}

// Controller defines the simple controller for player operations.
type Controller struct {
	repo playerRepository
}

// New creates a new instance of the player controller.
func New(r playerRepository) *Controller {
	return &Controller{repo: r}
}

// Get returns the player with the given id, or nil if no player with that id exists
func (c *Controller) Get(ctx context.Context, id int) (*model.Player, error) {
	return c.repo.Get(ctx, id)
}
