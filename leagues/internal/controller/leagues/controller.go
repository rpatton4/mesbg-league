package leagues

import (
	"context"
	"github.com/rpatton4/mesbg-league/leagues/pkg/model"
)

type leagueRepository interface {
	Get(ctx context.Context, id int) (*model.League, error)
}

// Controller defines the simple controller for league operations.
type Controller struct {
	repo leagueRepository
}

// New creates a new instance of the league controller.
func New(r leagueRepository) *Controller {
	return &Controller{repo: r}
}

// Get returns the league with the given id, or a svcerrors.NotFound if no league with that id exists
func (c *Controller) Get(ctx context.Context, id int) (*model.League, error) {
	return c.repo.Get(ctx, id)
}
