package memory

import (
	"context"
	"github.com/rpatton4/mesbg-league/player/pkg/model"
	"sync"
)

// Repository defines an in-memory repository for player data
type Repository struct {
	sync.RWMutex
	data map[string]*model.Person
}

func New() *Repository {
	return &Repository{data: map[string]*model.Person{}}
}

// Get retrieves a person by ID from the in-memory repository, if not person with the given
// ID exists, it returns nil.
func (r *Repository) Get(_ context.Context, id string) *model.Person {
	r.RLock()
	defer r.RUnlock()

	person, exists := r.data[id]
	if !exists {
		return nil
	}
	return person
}
