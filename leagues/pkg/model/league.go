package model

import (
	participants "github.com/rpatton4/mesbg-league/participants/pkg/model"
	rounds "github.com/rpatton4/mesbg-league/rounds/pkg/model"
)

// League is the model for a gaming league, which includes metadata about the league, its participants,
// and the games played within it.
type League struct {
	// ID is the unique identifier for the league
	ID string `json:"id"`

	// Active indicates whether the league is currently either not yet started or in progress (true) or has ended (false)
	Active bool `json:"active"`

	// Name of the league, e.g. "Fall 2025 Acme Gaming League"
	Name string `json:"name"`

	// Participants is a slice of participants (players + metadata) in the league
	Participants []*participants.Participant `json:"participants"`

	// Rounds is a slice of the rounds in the league, both those which have occurred and those which are upcoming
	Rounds []*rounds.Round `json:"rounds"`
	
	// NumberOfGames is the total number of games in the league
	NumberOfGames int `json:"numberOfGames"`

	// StartDate is the date of the first league game in YYYY-MM-DD format (ISO 8601)
	StartDate string `json:"startDate"`

	// EndDate is the date of the last league game in YYYY-MM-DD format (ISO 8601)
	EndDate string `json:"endDate"`

	// ExpectedDayOfWeek is the day of the week that games are generally expected to be played, e.g. "Monday", "Tuesday", etc.
	ExpectedDayOfWeek string `json:"expectedDayOfWeek"`
}
