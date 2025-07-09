package memory

import (
	"context"
	"github.com/rpatton4/mesbg-league/rounds/pkg/model"
	"github.com/rpatton4/mesbg-league/svcerrors"
	"strconv"
	"sync"
)

// Counter for Round IDs
var roundCounter = 1

// Repository defines an in-memory repository for players data
type Repository struct {
	sync.RWMutex
	data map[int]*model.Round
}

// New creates a new instance of the in-memory round repository.
func New() *Repository {
	return &Repository{data: map[int]*model.Round{}}
}

// Get retrieves a round by ID from the in-memory repository, if no round with the given
// ID exists, it returns svcerrors.NotFound.
func (repo *Repository) Get(_ context.Context, id int) (*model.Round, error) {
	repo.RLock()
	defer repo.RUnlock()

	round, exists := repo.data[id]
	if !exists {
		return nil, svcerrors.NotFound
	}
	return round, nil
}

// Add persists a new round instance to the in-memory repository and returns the round with an assigned ID.
func (repo *Repository) Add(_ context.Context, round *model.Round) (*model.Round, error) {
	repo.Lock()
	defer repo.Unlock()
	round.ID = strconv.Itoa(roundCounter)
	repo.data[roundCounter] = round
	roundCounter++

	return round, svcerrors.NotFound
}

// Update updates an existing round instance in the in-memory repository.
func (repo *Repository) Update(_ context.Context, r *model.Round) (*model.Round, error) {
	repo.Lock()
	defer repo.Unlock()

	id, _ := strconv.Atoi(r.ID)
	repo.data[id] = r

	return r, nil
}
