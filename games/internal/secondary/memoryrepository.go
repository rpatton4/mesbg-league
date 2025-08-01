package secondary

import (
	"context"
	"errors"
	"fmt"
	"github.com/rpatton4/mesbg-league/games/pkg"
	"github.com/rpatton4/mesbg-league/games/pkg/model"
	"github.com/rpatton4/mesbg-league/pkg/svcerrors"
	"strconv"
	"strings"
	"sync"
)

var gameCounter = 1

// MemoryRepository defines an in-memory repository (adapter) for the Games service
type MemoryRepository struct {
	sync.RWMutex
	data map[pkg.GameID]*model.Game
}

// NewMemoryRepository creates a new instance of the in-memory game repository.
func NewMemoryRepository() *MemoryRepository {
	return &MemoryRepository{data: map[pkg.GameID]*model.Game{}}
}

// GetByID retrieves a game by ID from the in-memory repository, if no game with the given
// ID exists, it returns ErrNotFound.
func (r *MemoryRepository) GetByID(_ context.Context, id pkg.GameID) (*model.Game, error) {
	r.RLock()
	defer r.RUnlock()

	g, exists := r.data[id]
	if !exists {
		return nil, svcerrors.ErrNotFound
	}

	return g, nil
}

// Create persists a new game instance to the in-memory repository and returns the game with an assigned ID.
func (r *MemoryRepository) Create(_ context.Context, g *model.Game) (*model.Game, error) {
	r.Lock()
	defer r.Unlock()

	if v, f, err := g.IsValid(); !v || err != nil {
		if !v || errors.Is(err, svcerrors.ErrModelInvalid) {
			return nil, fmt.Errorf("game %w "+strings.Join(f, "; "), svcerrors.ErrModelInvalid)
		} else if err != nil {
			return nil, err
		}

		return nil, fmt.Errorf("game %w: "+strings.Join(f, "; "), err)
	}

	g.ID = pkg.GameID(strconv.Itoa(gameCounter))
	r.data[g.ID] = g
	gameCounter++

	return g, nil
}

// Replace completely replaces an existing game instance with the provided one, using the ID from the provided game
// to find which game to replace. This cannot be used to create a new Game, and it is an idempotent operation.
// If the game is missing or invalid, this returns the appropriate svcerror
func (r *MemoryRepository) Replace(_ context.Context, g *model.Game) (*model.Game, error) {
	r.Lock()
	defer r.Unlock()

	if v, m, err := g.IsValid(); !v || err != nil {
		if !v || errors.Is(err, svcerrors.ErrModelInvalid) {
			return nil, fmt.Errorf("game %w: "+strings.Join(m, "; "), err)
		} else if err != nil {
			return nil, err
		}
	} else if g.ID == "" {
		return nil, fmt.Errorf("the game data sent with update is missing a game ID. Source: %w", svcerrors.ErrInvalidID)
	} else if r.data[g.ID] == nil {
		return nil, fmt.Errorf("the game with the given ID '%s' is not found. Source: %w", g.ID, svcerrors.ErrNotFound)
	}

	r.data[g.ID] = g
	return g, nil
}

// DeleteByID deletes an existing game instance in the in-memory repository. Returns true if the game was found and
// deleted, false otherwise. This is an idempotent operation.
func (r *MemoryRepository) DeleteByID(_ context.Context, id pkg.GameID) (bool, error) {
	r.Lock()
	defer r.Unlock()

	if r.data[id] != nil {
		r.data[id] = nil
		return true, nil
	}

	return false, svcerrors.ErrNotFound
}
