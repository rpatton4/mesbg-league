package secondary

import (
	games "github.com/rpatton4/mesbg-league/games/pkg"
	"github.com/rpatton4/mesbg-league/games/pkg/model"
	"testing"
)

// Sort of in passing this also tests that the memory repository meets the Repository interface spec
var r Repository

func TestMemoryRepoAddGame(t *testing.T) {
	r = NewMemoryRepository()
	g := createFakeGame()

	g, err := r.Create(nil, g)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if g.ID == games.GameID("") {
		t.Errorf("Expected game ID to be assigned, got 0")
	}
}

func TestMemoryRepoGetById(t *testing.T) {
	r = NewMemoryRepository()
	// Create the game to search for
	g := createFakeGame()

	g, err := r.Create(nil, g)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	if g.ID == games.GameID("") {
		t.Errorf("Expected game ID to be assigned, got 0")
	}

	// Search for it
	result, err := r.GetByID(nil, g.ID)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if result.ID != g.ID {
		t.Errorf("Expected game IDs to match, they didn't")
	}
}

func TestMemoryRepoReplaceGame(t *testing.T) {
	r = NewMemoryRepository()
	originalScore := 20
	updatedScore := 30
	originalRoundID := "999"

	// Create the game to search for
	g := createFakeGame()
	g.Side1TotalVictoryPoints = originalScore
	g.RoundID = originalRoundID

	g, err := r.Create(nil, g)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	if g.ID == games.GameID("") {
		t.Errorf("Expected game ID to be assigned, got 0")
	}

	g.Side1TotalVictoryPoints = updatedScore
	g.RoundID = "" // Clear out the round ID to put it back to unset

	g2, err2 := r.Replace(nil, g)
	if err2 != nil {
		t.Fatalf("Expected no error, got %v", err2)
	}
	if g2.Side1TotalVictoryPoints != updatedScore {
		t.Errorf("Expected score to be updated and it wasn't ")
	}
	if g2.RoundID != "" {
		t.Errorf("Expected RoundID to be un-set and it wasn't ")
	}

	g3, err3 := r.GetByID(nil, g2.ID)
	if err3 != nil {
		t.Fatalf("Expected no error, got %v", err3)
	}
	if g3.Side1TotalVictoryPoints != updatedScore {
		t.Errorf("Expected queried score to be updated and it wasn't ")
	}
	if g3.RoundID != "" {
		t.Errorf("Expected queried RoundID to be un-set and it wasn't ")
	}
}

func TestMemoryRepoDeleteGame(t *testing.T) {
	r = NewMemoryRepository()
	// Create the game to delete
	g := createFakeGame()
	g, err := r.Create(nil, g)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	if g.ID == games.GameID("") {
		t.Errorf("Expected game ID to be assigned, got 0")
	}

	success, err := r.DeleteByID(nil, g.ID)

	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if !success {
		t.Errorf("Expected delete to return true, got false")
	}

	// now delete it again, it should fail
	success, err = r.DeleteByID(nil, g.ID)
	if err == nil {
		t.Fatalf("Expected error while deleting non-existent game, but it succeeded")
	}

	if success {
		t.Errorf("Expected delete to return false, got true")
	}
}

func createFakeGame() *model.Game {
	return &model.Game{
		Side1ID:                 "123",
		Side2ID:                 "456",
		RoundID:                 "789",
		Side1TotalVictoryPoints: 10,
		Side2TotalVictoryPoints: 15,
		Side1KilledGeneral:      true,
		Side2KilledGeneral:      true,
		Status:                  model.GameStatePlayCompleted,
	}
}
