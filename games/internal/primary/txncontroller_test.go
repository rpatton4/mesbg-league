package primary

import (
	"github.com/rpatton4/mesbg-league/games/internal/secondary"
	games "github.com/rpatton4/mesbg-league/games/pkg"
	"github.com/rpatton4/mesbg-league/games/pkg/model"
	"testing"
)

func TestTxnControllerAddGame(t *testing.T) {
	ctrl := createController()
	g := createFakeGame()

	g, err := ctrl.Create(nil, g)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if g.ID == games.GameID("") {
		t.Errorf("Expected game ID to be assigned, got 0")
	}
}

func TestTxnControllerGetById(t *testing.T) {
	ctrl := createController()
	// Create the game to search for
	g := createFakeGame()

	g, err := ctrl.Create(nil, g)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	if g.ID == games.GameID("") {
		t.Errorf("Expected game ID to be assigned, got 0")
	}

	// Search for it
	result, err := ctrl.GetByID(nil, g.ID)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if result.ID != g.ID {
		t.Errorf("Expected game IDs to match, they didn't")
	}
}

func TestTxnControllerReplaceGame(t *testing.T) {
	ctrl := createController()
	originalScore := 20
	updatedScore := 30
	originalRoundID := "999"

	// Create the game to search for
	g := createFakeGame()
	g.Side1TotalVictoryPoints = originalScore
	g.RoundID = originalRoundID

	g, err := ctrl.Create(nil, g)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	if g.ID == games.GameID("") {
		t.Errorf("Expected game ID to be assigned, got 0")
	}

	g.Side1TotalVictoryPoints = updatedScore
	g.RoundID = "" // Clear out the round ID to put it back to unset

	g2, err2 := ctrl.Replace(nil, g)
	if err2 != nil {
		t.Fatalf("Expected no error, got %v", err2)
	}
	if g2.Side1TotalVictoryPoints != updatedScore {
		t.Errorf("Expected score to be updated and it wasn't ")
	}
	if g2.RoundID != "" {
		t.Errorf("Expected RoundID to be un-set and it wasn't ")
	}

	g3, err3 := ctrl.GetByID(nil, g2.ID)
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

func TestTxnControllerDeleteGame(t *testing.T) {
	ctrl := createController()
	// Create the game to delete
	g := createFakeGame()
	g, err := ctrl.Create(nil, g)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	if g.ID == games.GameID("") {
		t.Errorf("Expected game ID to be assigned, got 0")
	}

	success, err := ctrl.DeleteByID(nil, g.ID)

	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if !success {
		t.Errorf("Expected delete to return true, got false")
	}

	// now delete it again, it should fail
	success, err = ctrl.DeleteByID(nil, g.ID)
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
		Status:                  games.GameStatePlayCompleted,
	}
}

func createController() *TxnController {
	repo := secondary.NewMemoryRepository()
	return NewTxnController(repo)
}
