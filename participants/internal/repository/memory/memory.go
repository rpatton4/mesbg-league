package memory

import (
	"context"
	"github.com/rpatton4/mesbg-league/participants/pkg/model"
	"github.com/rpatton4/mesbg-league/svcerrors"
	"strconv"
	"sync"
)

var participantCounter = 1

// Repository defines an in-memory repository for participant data
type Repository struct {
	sync.RWMutex
	data map[int]*model.Participant
}

// New creates a new instance of the in-memory participant repository.
func New() *Repository {
	return &Repository{data: map[int]*model.Participant{}}
}

// Get retrieves a participant by ID from the in-memory repository, if no participant with the given
// ID exists, it returns NotFound.
func (r *Repository) Get(_ context.Context, id int) (*model.Participant, error) {
	r.RLock()
	defer r.RUnlock()

	p, exists := r.data[id]
	if !exists {
		return nil, svcerrors.NotFound
	}
	return p, nil
}

// Add persists a new participant instance to the in-memory repository and returns the participant with an assigned ID.
func (r *Repository) Add(_ context.Context, l *model.Participant) (*model.Participant, error) {
	r.Lock()
	defer r.Unlock()
	l.ID = strconv.Itoa(participantCounter)
	r.data[participantCounter] = l
	participantCounter++

	return l, svcerrors.NotFound
}

// Update updates an existing participant instance in the in-memory repository.
func (r *Repository) Update(_ context.Context, p *model.Participant) (*model.Participant, error) {
	r.Lock()
	defer r.Unlock()

	id, _ := strconv.Atoi(p.ID)
	r.data[id] = p

	return p, nil
}
