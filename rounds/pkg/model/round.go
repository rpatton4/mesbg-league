package model

import (
	gamesheader "github.com/rpatton4/mesbg-league/games/pkg/header"
	games "github.com/rpatton4/mesbg-league/games/pkg/model"
	leagues "github.com/rpatton4/mesbg-league/leagues/pkg/header"
	"github.com/rpatton4/mesbg-league/rounds/pkg/header"
)

// Round models one round of games in a league, linking the games scheduled and played for that round to the league
type Round struct {
	// ID is the unique identifier for the round
	ID header.RoundID `json:"id,omitempty"`

	// LeagueID is the key to the league this round belongs to
	LeagueID leagues.LeagueID `json:"leagueId"`

	// Number indicates which round this is in the league (1, 2, 3, etc)
	Number int `json:"number"`

	// ScenarioName is the name of the scenario expected to be played in this round, from the MSBG
	// rule book or the matched play guide
	ScenarioName string `json:"scenarioName"`

	// Games is the slice of games scheduled/played in this round
	Games []games.Game `json:"games,omitempty"`
}

type ShallowRound struct {
	// ID is the unique identifier for the round
	ID header.RoundID `json:"id,omitempty"`

	// LeagueID is the key to the league this round belongs to
	LeagueID leagues.LeagueID `json:"leagueId"`

	// Number indicates which round this is in the league (1, 2, 3, etc)
	Number int `json:"number"`

	// ScenarioName is the name of the scenario expected to be played in this round, from the MSBG
	// rule book or the matched play guide
	ScenarioName string `json:"scenarioName"`

	// GameIDs is the slice of IDs for games scheduled/played in this round
	GameIDs []gamesheader.GameID `json:"gameIDs,omitempty"`
}
