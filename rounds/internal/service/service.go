package service

import (
	"github.com/rpatton4/mesbg-league/rounds/pkg/model"
)

func RoundToLite(r *model.Round) *model.RoundLite {
	if r == nil {
		return nil
	}

	rl := &model.RoundLite{
		ID:           r.ID,
		LeagueID:     r.LeagueID,
		Number:       r.Number,
		ScenarioName: r.ScenarioName,
	}

	for _, g := range r.Games {
		rl.GameIDs = append(rl.GameIDs, g.ID)
	}

	return rl
}
