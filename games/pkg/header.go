package pkg

type GameID string

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
