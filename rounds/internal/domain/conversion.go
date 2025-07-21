package domain

import (
	"github.com/rpatton4/mesbg-league/rounds/pkg/model"
)

// DeepToShallow replaces any entities embedded in the Round with their IDs, returning a ShallowRound.
// The use case is meant to be for persisting a Round without including the data which is not owned by the Round,
func DeepToShallow(r *model.Round) *model.ShallowRound {
	if r == nil {
		return nil
	}

	sr := &model.ShallowRound{
		ID:           r.ID,
		LeagueID:     r.LeagueID,
		Number:       r.Number,
		ScenarioName: r.ScenarioName,
	}

	for _, g := range r.Games {
		sr.GameIDs = append(sr.GameIDs, g.ID)
	}

	return sr
}
