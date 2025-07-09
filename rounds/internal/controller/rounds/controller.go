package rounds

import (
	"context"
	"github.com/rpatton4/mesbg-league/rounds/pkg/model"
)

type roundRepository interface {
	Get(ctx context.Context, id int) (*model.Round, error)
}

// Controller defines the simple controller for round operations.
type Controller struct {
	repo roundRepository
}

// New creates a new instance of the round controller.
func New(repo roundRepository) *Controller {
	return &Controller{repo: repo}
}

// Get returns the round with the given id, or svcerrors.NotFound if no round with that id exists
func (c *Controller) Get(ctx context.Context, id int) (*model.Round, error) {
	return c.repo.Get(ctx, id)
}
