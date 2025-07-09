package participants

import (
	"context"
	"github.com/rpatton4/mesbg-league/participants/pkg/model"
)

type participantRepository interface {
	Get(ctx context.Context, id int) (*model.Participant, error)
}

// Controller defines the simple controller for participant operations.
type Controller struct {
	repo participantRepository
}

// New creates a new instance of the participant controller.
func New(r participantRepository) *Controller {
	return &Controller{repo: r}
}

// Get returns the participant with the given id, or a svcerrors.NotFound if no participant with that id exists
func (c *Controller) Get(ctx context.Context, id int) (*model.Participant, error) {
	return c.repo.Get(ctx, id)
}
