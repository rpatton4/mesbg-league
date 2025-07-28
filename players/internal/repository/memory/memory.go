package memory

import (
	"context"
	"github.com/rpatton4/mesbg-league/pkg/svcerrors"
	"github.com/rpatton4/mesbg-league/players/pkg/model"
	"strconv"
	"sync"
)

// Counter for Player IDs
var playerCounter = 1

// Repository defines an in-memory repository for players data
type Repository struct {
	sync.RWMutex
	data map[model.PlayerID]*model.Player
}

// New creates a new instance of the in-memory players repository.
func New() *Repository {
	return &Repository{data: map[model.PlayerID]*model.Player{}}
}

// GetByID retrieves a player by ID from the in-memory repository, if no player with the given
// ID exists, it returns ErrNotFound.
func (r *Repository) GetByID(_ context.Context, id model.PlayerID) (*model.Player, error) {
	r.RLock()
	defer r.RUnlock()

	p, exists := r.data[id]
	if !exists {
		return nil, svcerrors.ErrNotFound
	}
	return p, nil
}

// Create persists a new players instance to the in-memory repository and returns the players with an assigned ID.
func (r *Repository) Create(_ context.Context, p *model.Player) (*model.Player, error) {
	r.Lock()
	defer r.Unlock()

	p.ID = model.PlayerID(strconv.Itoa(playerCounter))
	r.data[p.ID] = p
	playerCounter++

	return p, nil
}

// Replace completely replaces an existing player instance with the provided one, using the ID from the provided player
// to find which player to replace. This cannot be used to create a new player, and it is an idempotent operation.
// This is an intended equivalent to the HTTP PUT operation, though it purposefully does not allow the create which
// PUT is sometimes interpreted as allowing (because that leaves ID creation up to the client).
func (r *Repository) Replace(_ context.Context, p *model.Player) (*model.Player, error) {
	r.Lock()
	defer r.Unlock()

	if p.ID == "" || r.data[p.ID] == nil {
		return nil, svcerrors.ErrInvalidID
	}
	r.data[p.ID] = p
	return p, nil
}

// DeleteByID deletes an existing player instance in the in-memory repository. Returns true if the player was found and
// deleted, false otherwise. This is an idempotent operation.
func (r *Repository) DeleteByID(_ context.Context, id model.PlayerID) bool {
	r.Lock()
	defer r.Unlock()

	if r.data[id] != nil {
		r.data[id] = nil
		return true
	}

	return false
}
