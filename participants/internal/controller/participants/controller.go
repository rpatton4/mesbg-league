package participants

import (
	"context"
	"errors"
	"github.com/rpatton4/mesbg-league/participants/pkg/model"
)

type participantRepository interface {
	GetByID(ctx context.Context, id model.ParticipantID) (*model.Participant, error)
	Create(ctx context.Context, p *model.Participant) (*model.Participant, error)
	Replace(ctx context.Context, p *model.Participant) (*model.Participant, error)
	DeleteByID(ctx context.Context, id model.ParticipantID) bool
}

// Controller defines the simple controller for participant operations.
type Controller struct {
	repo participantRepository
}

// New creates a new instance of the participant controller.
func New(r participantRepository) *Controller {
	return &Controller{repo: r}
}

// GetByID returns the participant with the given id, or svcerrors.NotFound if no participant with that id exists
func (c *Controller) GetByID(ctx context.Context, id model.ParticipantID) (*model.Participant, error) {
	return c.repo.GetByID(ctx, id)
}

// Create persists a new participant instance to the repository and returns the participant with an assigned ID.
// A generic error is returned if the participant to created is missing, while specific validation errors are
// passed along from the repository if the participant is invalid in some way.
func (c *Controller) Create(ctx context.Context, p *model.Participant) (*model.Participant, error) {
	if p == nil {
		return nil, errors.New("the participant to be created cannot be nil")
	}
	return c.repo.Create(ctx, p)
}

// Replace updates an existing participant in the repository with the provided participant.
// A generic error is returned if the participant to replaced is not present in the data store.
func (c *Controller) Replace(ctx context.Context, p *model.Participant) (*model.Participant, error) {
	if p == nil {
		return nil, errors.New("the participant to be created cannot be nil")
	}
	return c.repo.Replace(ctx, p)
}

// DeleteByID removes the participant with the given id from the repository. Returns true if the participant was found and
// deleted, false otherwise. This is an idempotent operation.
func (c *Controller) DeleteByID(ctx context.Context, id model.ParticipantID) bool {
	if id == "" {
		return false
	}
	return c.repo.DeleteByID(ctx, id)
}
