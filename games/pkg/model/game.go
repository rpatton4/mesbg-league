package model

import "github.com/rpatton4/mesbg-league/games/pkg/header"

// GameState is use to indicate the state of play for a game
type GameState int

const (
	GameStateNotStarted    GameState = 0
	GameStateInProgress    GameState = 1
	GameStatePlayCompleted GameState = 2
	GameStateBye           GameState = 3
	GameStateConceded      GameState = 4
	GameStateCancelled     GameState = 5
)

type Game struct {
	// ID is the unique identifier for the game
	ID header.GameID `json:"id,omitempty"`

	// Side1ID is the identifier of the player for the first side in the game.
	// First does not imply that this player acts first, it is simply a designator
	Side1ID string `json:"side1Id"`

	// Side2ID is the identifier of the player for the second side in the game.
	// Second does not imply that this player acts second, it is simply a designator
	Side2ID string `json:"side2Id"`

	// RoundID is the key to the round in a league which this game is part of
	RoundID string `json:"roundId"`

	// Side1TotalVictoryPoints is the total victory points scored by the first side in the game
	Side1TotalVictoryPoints int `json:"side1TotalVictoryPoints,omitempty"`

	// Side2TotalVictoryPoints is the total victory points scored by the second side in the game
	Side2TotalVictoryPoints int `json:"side2TotalVictoryPoints,omitempty"`

	// Side1TotalGeneralsKilled is true if the side 1 player killed the opposing general
	Side1KilledGeneral bool `json:"side1KilledGeneral,omitempty"`

	// Side2TotalGeneralsKilled is true if the side 2 player killed the opposing general
	Side2KilledGeneral bool `json:"side2KilledGeneral,omitempty"`

	// Status is used to track whether the game is scheduled, played, conceded etc.
	// See the GameStateXYZ constants for potential values.
	Status GameState `json:"status,omitempty"`
}
