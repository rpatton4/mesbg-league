package memory

import (
	"context"
	"github.com/rpatton4/mesbg-league/games/pkg/model"
	"github.com/rpatton4/mesbg-league/svcerrors"
	"strconv"
	"sync"
)

var gameCounter = 1

// Repository defines an in-memory repository for game data
type Repository struct {
	sync.RWMutex
	data map[model.GameID]*model.Game
}

// New creates a new instance of the in-memory game repository.
func New() *Repository {
	return &Repository{data: map[model.GameID]*model.Game{}}
}

// GetByID retrieves a game by ID from the in-memory repository, if no game with the given
// ID exists, it returns NotFound.
func (r *Repository) GetByID(_ context.Context, id model.GameID) (*model.Game, error) {
	r.RLock()
	defer r.RUnlock()

	g, exists := r.data[id]
	if !exists {
		return nil, svcerrors.NotFound
	}
	return g, nil
}

// Create persists a new game instance to the in-memory repository and returns the game with an assigned ID.
func (r *Repository) Create(_ context.Context, g *model.Game) (*model.Game, error) {
	r.Lock()
	defer r.Unlock()

	g.ID = model.GameID(strconv.Itoa(gameCounter))
	r.data[g.ID] = g
	gameCounter++

	return g, nil
}

// Replace completely replaces an existing game instance with the provided one, using the ID from the provided game
// to find which game to replace. This cannot be used to create a new Game, and it is an idempotent operation.
// This is an intended equivalent to the HTTP PUT operation, though it purposefully does not allow the create which
// PUT is sometimes interpreted as allowing (because that leaves ID creation up to the client).
func (r *Repository) Replace(_ context.Context, g *model.Game) (*model.Game, error) {
	r.Lock()
	defer r.Unlock()

	if g.ID == "" || r.data[g.ID] == nil {
		return nil, svcerrors.InvalidID
	}
	r.data[g.ID] = g
	return g, nil
}

// Delete deletes an existing game instance in the in-memory repository.
func (r *Repository) Delete(_ context.Context, id model.GameID) error {
	r.Lock()
	defer r.Unlock()

	r.data[id] = nil

	return nil
}
