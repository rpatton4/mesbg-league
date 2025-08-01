package primary

import (
	"context"
	"errors"
	"fmt"
	"github.com/danielgtaylor/huma/v2"
	mock_primary "github.com/rpatton4/mesbg-league/games/internal/primary/mocks"
	games "github.com/rpatton4/mesbg-league/games/pkg"
	"github.com/rpatton4/mesbg-league/games/pkg/model"
	"github.com/rpatton4/mesbg-league/pkg/svcerrors"
	players "github.com/rpatton4/mesbg-league/players/pkg"
	rounds "github.com/rpatton4/mesbg-league/rounds/pkg"
	"go.uber.org/mock/gomock"
	"testing"
)

func TestHumaHandlerMockedPut(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	mockController := mock_primary.NewMockSingleController(mockCtrl)
	handler := NewHumaHandler(mockController)
	var statusError huma.StatusError

	validGame := model.Game{
		ID:      games.GameID("1"),
		Side1ID: players.PlayerID("8"),
		Side2ID: players.PlayerID("9"),
		RoundID: rounds.RoundID("7"),
		Status:  games.GameStateNotStarted,
	}

	updatedGame := model.Game{
		ID:      games.GameID("1"),
		Side1ID: players.PlayerID("8"),
		Side2ID: players.PlayerID("9"),
		RoundID: rounds.RoundID("7"),
		Status:  games.GameStateNotStarted,
	}

	updatedGameMissingID := model.Game{
		Side1ID: players.PlayerID("8"),
		Side2ID: players.PlayerID("9"),
		RoundID: rounds.RoundID("7"),
		Status:  games.GameStateNotStarted,
	}

	invalidGame := model.Game{
		ID:      games.GameID("3"),
		Side1ID: players.PlayerID("8"),
		RoundID: rounds.RoundID("7"),
	}

	notFoundGame := model.Game{
		ID:      games.GameID("999"),
		Side1ID: players.PlayerID("8"),
		Side2ID: players.PlayerID("9"),
		RoundID: rounds.RoundID("7"),
		Status:  games.GameStateNotStarted,
	}

	// Prepare the mock
	mockController.EXPECT().Replace(gomock.Any(), &validGame).Return(&updatedGame, nil).Times(1)                                                                                   // valid, exists
	mockController.EXPECT().Replace(gomock.Any(), &invalidGame).Return(nil, fmt.Errorf("game %w:", svcerrors.ErrModelInvalid)).Times(1)                                            // invalid
	mockController.EXPECT().Replace(gomock.Any(), &notFoundGame).Return(nil, fmt.Errorf("the game with the given ID '%s' is %w", notFoundGame.ID, svcerrors.ErrNotFound)).Times(1) // valid, not found
	mockController.EXPECT().Replace(gomock.Any(), nil).Return(nil, svcerrors.ErrModelMissing).Times(1)
	mockController.EXPECT().Replace(gomock.Any(), &updatedGameMissingID).Return(nil, fmt.Errorf("the game data sent with update is missing a game ID. Source: %w", svcerrors.ErrInvalidID)) // trying to update with a missing ID

	// Test for valid game update
	res, err := handler.Put(context.Background(), &PutRequest{Body: &validGame})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if res.Body.ID != updatedGame.ID {
		t.Fatalf("expected game ID '%s', got '%s'", validGame.ID, res.Body.ID)
	}

	// Test for invalid game when trying to update
	_, err = handler.Put(context.Background(), &PutRequest{Body: &invalidGame})
	if err == nil {
		t.Fatalf("expected error for invalid game, got nil")
	} else if errors.As(err, &statusError) && statusError.GetStatus() != 400 {
		t.Fatalf("expected 400 for invalid game, got %v", err)
	}

	// Test for game which can't be found when trying to update
	_, err = handler.Put(context.Background(), &PutRequest{Body: &notFoundGame})
	if err == nil {
		t.Fatalf("expected error for unfound game, got nil")
	} else if errors.As(err, &statusError) && statusError.GetStatus() != 400 {
		t.Fatalf("expected 400 for unfound game, got %v", err)
	}

	// Test for missing game when trying to update
	_, err = handler.Put(context.Background(), &PutRequest{Body: nil})
	if err == nil {
		t.Fatalf("expected error for empty game, got nil")
	} else if errors.As(err, &statusError) && statusError.GetStatus() != 400 {
		t.Fatalf("expected 400 for empty game, got %v", err)
	}

	// Test for missing game when trying to update
	_, err = handler.Put(context.Background(), &PutRequest{Body: &updatedGameMissingID})
	if err == nil {
		t.Fatalf("expected error for game without ID, got nil")
	} else if errors.As(err, &statusError) && statusError.GetStatus() != 400 {
		t.Fatalf("expected 400 for game without ID, got %v", err)
	}
}

func TestHumaHandlerMockedPost(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	mockController := mock_primary.NewMockSingleController(mockCtrl)
	handler := NewHumaHandler(mockController)
	var statusError huma.StatusError
	validGame := model.Game{
		Side1ID: players.PlayerID("8"),
		Side2ID: players.PlayerID("9"),
		RoundID: rounds.RoundID("7"),
		Status:  games.GameStateNotStarted,
	}

	invalidGame := model.Game{
		Side1ID: players.PlayerID("8"),
		RoundID: rounds.RoundID("7"),
	}

	game1Return := model.Game{
		ID:      games.GameID("1"),
		Side1ID: players.PlayerID("8"),
		Side2ID: players.PlayerID("9"),
		RoundID: rounds.RoundID("7"),
		Status:  games.GameStateNotStarted,
	}

	// Prepare the mock
	mockController.EXPECT().Create(gomock.Any(), &validGame).Return(&game1Return, nil).Times(1)
	mockController.EXPECT().Create(gomock.Any(), &invalidGame).Return(nil, fmt.Errorf("game %w: Side2ID=''; Status is not set", svcerrors.ErrModelInvalid)).Times(1)
	mockController.EXPECT().Create(gomock.Any(), nil).Return(nil, svcerrors.ErrModelMissing).Times(1)

	// Test for valid game creation
	res, err := handler.Post(context.Background(), &PostRequest{Body: &validGame})

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if res.Body.ID != game1Return.ID {
		t.Fatalf("expected game ID '%s', got '%s'", validGame.ID, res.Body.ID)
	}

	// Test for missing game when trying to create
	_, err = handler.Post(context.Background(), &PostRequest{Body: nil})
	if err == nil {
		t.Fatalf("expected error for empty game, got nil")
	} else if errors.As(err, &statusError) && statusError.GetStatus() != 400 {
		t.Fatalf("expected 400 for empty game, got %v", err)
	}

	// Test for invalid game when trying to create
	_, err = handler.Post(context.Background(), &PostRequest{Body: &invalidGame})
	if err == nil {
		t.Fatalf("expected error for invalid game, got nil")
	} else if errors.As(err, &statusError) && statusError.GetStatus() != 400 {
		t.Fatalf("expected 400 for invalid game, got %v", err)
	}
}

func TestHumaHandlerMockedDelete(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	mockController := mock_primary.NewMockSingleController(mockCtrl)
	handler := NewHumaHandler(mockController)
	var statusError huma.StatusError

	// Prepare the mock
	mockController.EXPECT().DeleteByID(gomock.Any(), games.GameID("1")).Return(true, nil).Times(1)
	mockController.EXPECT().DeleteByID(gomock.Any(), games.GameID("999")).Return(false, svcerrors.ErrNotFound).Times(1) // not found

	// Test for valid game deletion
	_, err := handler.Delete(context.Background(), &DeleteRequest{ID: games.GameID("1")})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	// Test for game which can't be found when trying to delete
	_, err = handler.Delete(context.Background(), &DeleteRequest{ID: games.GameID("999")})
	if err == nil {
		t.Fatalf("expected error for unfound game, got nil")
	}
	if !errors.As(err, &statusError) || statusError.GetStatus() != 404 {
		t.Fatalf("expected 404 for unfound game, got %v", err)
	}
}

func TestHumaHandlerMockedGetByID(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	mockController := mock_primary.NewMockSingleController(mockCtrl)
	handler := NewHumaHandler(mockController)
	var statusError huma.StatusError
	game1 := model.Game{
		ID:      games.GameID("1"),
		Side1ID: players.PlayerID("8"),
		Side2ID: players.PlayerID("9"),
		RoundID: rounds.RoundID("7"),
		Status:  games.GameStateNotStarted,
	}

	// Prepare the mock
	mockController.EXPECT().GetByID(gomock.Any(), games.GameID("1")).Return(&game1, nil).Times(1)                  // valid, exists
	mockController.EXPECT().GetByID(gomock.Any(), games.GameID("999")).Return(nil, svcerrors.ErrNotFound).Times(1) // valid, not found
	mockController.EXPECT().GetByID(gomock.Any(), games.GameID("")).Return(nil, svcerrors.ErrInvalidID).Times(1)   // invalid

	// Test valid ID
	res, err := handler.GetByID(context.Background(), &GetByIDRequest{ID: games.GameID("1")})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if res.Body.ID != game1.ID {
		t.Fatalf("expected game ID '%s', got '%s'", game1.ID, res.Body.ID)
	}

	// Test ID which is unknown
	_, err = handler.GetByID(context.Background(), &GetByIDRequest{ID: games.GameID("999")})
	if err == nil {
		t.Fatalf("expected error for non-existent game, got nil")
	} else if errors.As(err, &statusError) && statusError.GetStatus() != 404 {
		t.Fatalf("expected 404 for non-existent game, got %v", err)
	} else if !errors.As(err, &statusError) {
		t.Fatalf("Incorrect error while searching for game by ID with an invalid ID, got %v", err)
	}

	// Test missing ID
	_, err = handler.GetByID(context.Background(), &GetByIDRequest{ID: games.GameID("")})
	if err == nil {
		t.Fatalf("expected error for empty game ID, got nil")
	} else if errors.As(err, &statusError) && statusError.GetStatus() != 400 {
		t.Fatalf("expected 400 for empty game ID, got %v", err)
	}
}
