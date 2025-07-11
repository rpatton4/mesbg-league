package memory

import (
	"context"
	"github.com/rpatton4/mesbg-league/players/pkg/model"
	"github.com/rpatton4/mesbg-league/svcerrors"
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

// Get retrieves a person by ID from the in-memory repository, if no person with the given
// ID exists, it returns nil.
func (r *Repository) Get(_ context.Context, id model.PlayerID) (*model.Player, error) {
	r.RLock()
	defer r.RUnlock()

	person, exists := r.data[id]
	if !exists {
		return nil, svcerrors.NotFound
	}
	return person, nil
}

// Add persists a new players instance to the in-memory repository and returns the players with an assigned ID.
func (r *Repository) Add(_ context.Context, player *model.Player) (*model.Player, error) {
	r.Lock()
	defer r.Unlock()
	player.ID = model.PlayerID(strconv.Itoa(playerCounter))
	r.data[player.ID] = player
	playerCounter++

	return player, svcerrors.NotFound
}

// Update updates an existing players instance in the in-memory repository.
func (r *Repository) Update(_ context.Context, p *model.Player) (*model.Player, error) {
	r.Lock()
	defer r.Unlock()

	r.data[p.ID] = p

	return p, nil
}
