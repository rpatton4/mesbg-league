package players

import (
	"context"
	"errors"
	players "github.com/rpatton4/mesbg-league/players/pkg"
	"github.com/rpatton4/mesbg-league/players/pkg/model"
)

type playerRepository interface {
	GetByID(ctx context.Context, id players.PlayerID) (*model.Player, error)
	Create(ctx context.Context, p *model.Player) (*model.Player, error)
	Replace(ctx context.Context, p *model.Player) (*model.Player, error)
	DeleteByID(ctx context.Context, id players.PlayerID) bool
}

// Controller defines the simple controller for player operations.
type Controller struct {
	repo playerRepository
}

// New creates a new instance of the players controller.
func New(r playerRepository) *Controller {
	return &Controller{repo: r}
}

// GetByID returns the player with the given id, or svcerrors.NotFound if no player with that id exists
func (c *Controller) GetByID(ctx context.Context, id players.PlayerID) (*model.Player, error) {
	return c.repo.GetByID(ctx, id)
}

// Create persists a new player instance to the repository and returns the player with an assigned ID.
// A generic error is returned if the player to created is missing, while specific validation errors are
// passed along from the repository if the player is invalid in some way.
func (c *Controller) Create(ctx context.Context, p *model.Player) (*model.Player, error) {
	if p == nil {
		return nil, errors.New("the player to be created cannot be nil")
	}
	return c.repo.Create(ctx, p)
}

// Replace updates an existing player in the repository with the provided player.
// A generic error is returned if the player to replaced is not present in the data store.
func (c *Controller) Replace(ctx context.Context, p *model.Player) (*model.Player, error) {
	if p == nil {
		return nil, errors.New("the player to be created cannot be nil")
	}
	return c.repo.Replace(ctx, p)
}

// DeleteByID removes the player with the given id from the repository. Returns true if the player was found and
// deleted, false otherwise. This is an idempotent operation.
func (c *Controller) DeleteByID(ctx context.Context, id players.PlayerID) bool {
	if id == "" {
		return false
	}
	return c.repo.DeleteByID(ctx, id)
}
