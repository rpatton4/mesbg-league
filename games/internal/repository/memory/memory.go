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
	data map[int]*model.Game
}

// New creates a new instance of the in-memory game repository.
func New() *Repository {
	return &Repository{data: map[int]*model.Game{}}
}

// Get retrieves a game by ID from the in-memory repository, if no game with the given
// ID exists, it returns NotFound.
func (r *Repository) Get(_ context.Context, id int) (*model.Game, error) {
	r.RLock()
	defer r.RUnlock()

	p, exists := r.data[id]
	if !exists {
		return nil, svcerrors.NotFound
	}
	return p, nil
}

// Add persists a new game instance to the in-memory repository and returns the game with an assigned ID.
func (r *Repository) Add(_ context.Context, l *model.Game) (*model.Game, error) {
	r.Lock()
	defer r.Unlock()
	l.ID = strconv.Itoa(gameCounter)
	r.data[gameCounter] = l
	gameCounter++

	return l, svcerrors.NotFound
}

// Update updates an existing participant instance in the in-memory repository.
func (r *Repository) Update(_ context.Context, p *model.Game) (*model.Game, error) {
	r.Lock()
	defer r.Unlock()

	id, _ := strconv.Atoi(p.ID)
	r.data[id] = p

	return p, nil
}
