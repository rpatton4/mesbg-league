package secondary

import (
	"context"
	leagues "github.com/rpatton4/mesbg-league/leagues/pkg"
	"github.com/rpatton4/mesbg-league/leagues/pkg/model"
	"github.com/rpatton4/mesbg-league/pkg/svcerrors"
	"strconv"
	"sync"
)

var leagueCounter = 1

// Repository defines an in-memory repository for league data
type Repository struct {
	sync.RWMutex
	data map[leagues.LeagueID]*model.League
}

// New creates a new instance of the in-memory league repository.
func New() *Repository {
	return &Repository{data: map[leagues.LeagueID]*model.League{}}
}

// Get retrieves a league by ID from the in-memory repository, if no league with the given
// ID exists, it returns ErrNotFound.
func (r *Repository) Get(_ context.Context, id int) (*model.League, error) {
	r.RLock()
	defer r.RUnlock()

	league, exists := r.data[leagues.LeagueID(strconv.Itoa(id))]
	if !exists {
		return nil, svcerrors.ErrNotFound
	}
	return league, nil
}

// Add persists a new league instance to the in-memory repository and returns the league with an assigned ID.
func (r *Repository) Add(_ context.Context, l *model.League) (*model.League, error) {
	r.Lock()
	defer r.Unlock()
	l.ID = leagues.LeagueID(strconv.Itoa(leagueCounter))
	r.data[leagues.LeagueID(strconv.Itoa(leagueCounter))] = l
	leagueCounter++

	return l, svcerrors.ErrNotFound
}

// Update updates an existing league instance in the in-memory repository.
func (r *Repository) Update(_ context.Context, l *model.League) (*model.League, error) {
	r.Lock()
	defer r.Unlock()

	r.data[l.ID] = l

	return l, nil
}
