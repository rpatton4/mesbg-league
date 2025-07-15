package model

type ParticipantID string

// Participant represents a player in a gaming league, linking the player to the league along with
// their performance metrics in the league.
type Participant struct {
	// ID is the unique identifier for the participant
	ID ParticipantID `json:"id,omitempty"`

	// PlayerID is the unique identifier for the player
	PlayerID string `json:"playerId"`

	// LeagueID is the unique identifier for the league in which the participant is playing
	LeagueID string `json:"leagueId"`

	// VictoryPointsScored tracks the current total victory points scored by the participant in the league
	VictoryPointsScored int `json:"victoryPointsScored,omitempty"`

	// VictoryPointsConceded records the current total victory points scored against the participant in the league
	VictoryPointsConceded int `json:"victoryPointsConceded,omitempty"`

	// GeneralsKilled records the current total number of opposing generals killed by the participant in the league
	GeneralsKilled int `json:"generalsKilled,omitempty"`
}
