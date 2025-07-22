package outbound

import (
	"context"
	"github.com/rpatton4/mesbg-league/games/pkg/header"
	"github.com/rpatton4/mesbg-league/games/pkg/model"
	"github.com/rpatton4/mesbg-league/pkg/svcerrors"
	"strconv"
	"sync"
)

var gameCounter = 1

// MemoryRepository defines an in-memory repository (adapter) for the Games service
type MemoryRepository struct {
	sync.RWMutex
	data map[header.GameID]*model.Game
}

// NewMemoryRepository creates a new instance of the in-memory game repository.
func NewMemoryRepository() *MemoryRepository {
	return &MemoryRepository{data: map[header.GameID]*model.Game{}}
}

// GetByID retrieves a game by ID from the in-memory repository, if no game with the given
// ID exists, it returns NotFound.
func (r *MemoryRepository) GetByID(_ context.Context, id header.GameID) (*model.Game, error) {
	r.RLock()
	defer r.RUnlock()

	g, exists := r.data[id]
	if !exists {
		return nil, svcerrors.NotFound
	}

	return g, nil
}

// Create persists a new game instance to the in-memory repository and returns the game with an assigned ID.
func (r *MemoryRepository) Create(_ context.Context, g *model.Game) (*model.Game, error) {
	r.Lock()
	defer r.Unlock()

	g.ID = header.GameID(strconv.Itoa(gameCounter))
	r.data[g.ID] = g
	gameCounter++

	return g, nil
}

// Replace completely replaces an existing game instance with the provided one, using the ID from the provided game
// to find which game to replace. This cannot be used to create a new Game, and it is an idempotent operation.
// This is an intended equivalent to the HTTP PUT operation, though it purposefully does not allow the create which
// PUT is sometimes interpreted as allowing (because that leaves ID creation up to the client).
func (r *MemoryRepository) Replace(_ context.Context, g *model.Game) (*model.Game, error) {
	r.Lock()
	defer r.Unlock()

	if g.ID == "" || r.data[g.ID] == nil {
		return nil, svcerrors.InvalidID
	}
	r.data[g.ID] = g
	return g, nil
}

// DeleteByID deletes an existing game instance in the in-memory repository. Returns true if the game was found and
// deleted, false otherwise. This is an idempotent operation.
func (r *MemoryRepository) DeleteByID(_ context.Context, id header.GameID) bool {
	r.Lock()
	defer r.Unlock()

	if r.data[id] != nil {
		r.data[id] = nil
		return true
	}

	return false
}
