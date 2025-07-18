package memory

import (
	"context"
	"github.com/rpatton4/mesbg-league/leagues/pkg/model"
	"github.com/rpatton4/mesbg-league/pkg/svcerrors"
	"strconv"
	"sync"
)

var leagueCounter = 1

// Repository defines an in-memory repository for league data
type Repository struct {
	sync.RWMutex
	data map[int]*model.League
}

// New creates a new instance of the in-memory league repository.
func New() *Repository {
	return &Repository{data: map[int]*model.League{}}
}

// Get retrieves a league by ID from the in-memory repository, if no league with the given
// ID exists, it returns NotFound.
func (r *Repository) Get(_ context.Context, id int) (*model.League, error) {
	r.RLock()
	defer r.RUnlock()

	league, exists := r.data[id]
	if !exists {
		return nil, svcerrors.NotFound
	}
	return league, nil
}

// Add persists a new league instance to the in-memory repository and returns the league with an assigned ID.
func (r *Repository) Add(_ context.Context, l *model.League) (*model.League, error) {
	r.Lock()
	defer r.Unlock()
	l.ID = strconv.Itoa(leagueCounter)
	r.data[leagueCounter] = l
	leagueCounter++

	return l, svcerrors.NotFound
}

// Update updates an existing league instance in the in-memory repository.
func (r *Repository) Update(_ context.Context, l *model.League) (*model.League, error) {
	r.Lock()
	defer r.Unlock()

	id, _ := strconv.Atoi(l.ID)
	r.data[id] = l

	return l, nil
}
