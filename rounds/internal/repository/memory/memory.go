package memory

import (
	"context"
	"github.com/rpatton4/mesbg-league/pkg/svcerrors"
	rounds "github.com/rpatton4/mesbg-league/rounds/pkg"
	"github.com/rpatton4/mesbg-league/rounds/pkg/model"
	"strconv"
	"sync"
)

// Counter for Round IDs
var roundCounter = 1

// Repository defines an in-memory repository for players data
type Repository struct {
	sync.RWMutex
	data map[rounds.RoundID]*model.Round
}

// New creates a new instance of the in-memory round repository.
func New() *Repository {
	return &Repository{data: map[rounds.RoundID]*model.Round{}}
}

// Get retrieves a round by ID from the in-memory repository, if no round with the given
// ID exists, it returns svcerrors.NotFound.
func (repo *Repository) Get(_ context.Context, id int) (*model.Round, error) {
	repo.RLock()
	defer repo.RUnlock()

	round, exists := repo.data[rounds.RoundID(strconv.Itoa(id))]
	if !exists {
		return nil, svcerrors.ErrNotFound
	}
	return round, nil
}

// Add persists a new round instance to the in-memory repository and returns the round with an assigned ID.
func (repo *Repository) Add(_ context.Context, round *model.Round) (*model.Round, error) {
	repo.Lock()
	defer repo.Unlock()
	round.ID = rounds.RoundID(strconv.Itoa(roundCounter))
	repo.data[rounds.RoundID(strconv.Itoa(roundCounter))] = round
	roundCounter++

	return round, svcerrors.ErrNotFound
}

// Update updates an existing round instance in the in-memory repository.
func (repo *Repository) Update(_ context.Context, r *model.Round) (*model.Round, error) {
	repo.Lock()
	defer repo.Unlock()

	repo.data[r.ID] = r

	return r, nil
}
