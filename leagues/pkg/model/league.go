package model

import "github.com/rpatton4/mesbg-league/players/pkg/model"

type League struct {
	ID                string `json:"id"`
	Active            bool   `json:"active"`
	Name              string `json:"name"`
	Participant       []*model.Player
	NumberOfGames     int    `json:"numberOfGames"`
	StartDate         string `json:"startDate"`
	EndDate           string `json:"endDate"`
	ExpectedDayOfWeek string `json:"expectedDayOfWeek"`
}
