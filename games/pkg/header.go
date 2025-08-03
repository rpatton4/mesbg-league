// Package pkg holds the common types and constants used across the games service. Model definitions are
// purposefully separate in order to avoid cyclical dependencies with other services' models.
package pkg

// GameID is a universally unique identifier for a game and is guaranteed to be safe for use in a URI
type GameID string

// GameState is use to indicate the state of play for a game, from not started to completed
type GameState int

const (
	// GameStateNotStarted indicates the game has been initialized with the required fields but has not been played yet, basically it is planned
	GameStateNotStarted GameState = 0

	// GameStateInProgress indicates the game is currently being played, and is used in the rare case where a game is not completed in one session or some
	// sort of competitive dashboard is being used to show the current state of play
	GameStateInProgress GameState = 1

	// GameStatePlayCompleted indicates the game has been played and the results have been recorded
	GameStatePlayCompleted GameState = 2

	// GameStateBye indicates the game was not played, but a player received a bye for some reason, such as there being
	// an odd number of players in a tournament or league
	GameStateBye GameState = 3

	// GameStateConceded indicates the game was not played, but one player conceded to another
	GameStateConceded GameState = 4

	// GameStateCancelled indicates the game was not played, but was cancelled for some reason, such as a player dropping out of a tournament or league
	GameStateCancelled GameState = 5
)
