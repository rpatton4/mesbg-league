package memory

import (
	"context"
	"github.com/rpatton4/mesbg-league/participants/pkg/model"
	"github.com/rpatton4/mesbg-league/pkg/svcerrors"
	"strconv"
	"sync"
)

var participantCounter = 1

// Repository defines an in-memory repository for participant data
type Repository struct {
	sync.RWMutex
	data map[model.ParticipantID]*model.Participant
}

// New creates a new instance of the in-memory participant repository.
func New() *Repository {
	return &Repository{data: map[model.ParticipantID]*model.Participant{}}
}

// Get retrieves a participant by ID from the in-memory repository, if no participant with the given
// ID exists, it returns NotFound.
func (r *Repository) GetByID(_ context.Context, id model.ParticipantID) (*model.Participant, error) {
	r.RLock()
	defer r.RUnlock()

	p, exists := r.data[id]
	if !exists {
		return nil, svcerrors.NotFound
	}
	return p, nil
}

// Add persists a new participant instance to the in-memory repository and returns the participant with an assigned ID.
func (r *Repository) Create(_ context.Context, p *model.Participant) (*model.Participant, error) {
	r.Lock()
	defer r.Unlock()
	p.ID = model.ParticipantID(strconv.Itoa(participantCounter))
	r.data[p.ID] = p
	participantCounter++

	return p, svcerrors.NotFound
}

// Replace completely replaces an existing participant instance with the provided one, using the ID from the provided participant
// to find which participant to replace. This cannot be used to create a new participant, and it is an idempotent operation.
// This is an intended equivalent to the HTTP PUT operation, though it purposefully does not allow the create which
// PUT is sometimes interpreted as allowing (because that leaves ID creation up to the client).
func (r *Repository) Replace(_ context.Context, p *model.Participant) (*model.Participant, error) {
	r.Lock()
	defer r.Unlock()

	if p.ID == "" || r.data[p.ID] == nil {
		return nil, svcerrors.InvalidID
	}
	r.data[p.ID] = p
	return p, nil
}

// DeleteByID deletes an existing participant instance in the in-memory repository. Returns true if the participant was found and
// deleted, false otherwise. This is an idempotent operation.
func (r *Repository) DeleteByID(_ context.Context, id model.ParticipantID) bool {
	r.Lock()
	defer r.Unlock()

	if r.data[id] != nil {
		r.data[id] = nil
		return true
	}

	return false
}
