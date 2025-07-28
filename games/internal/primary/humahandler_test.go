package primary

import (
	"context"
	"errors"
	games "github.com/rpatton4/mesbg-league/games/pkg"
	"github.com/rpatton4/mesbg-league/games/pkg/model"
	"github.com/rpatton4/mesbg-league/pkg/svcerrors"
	"testing"
)

type MockSingleController struct {
}

func (m *MockSingleController) GetByID(ctx context.Context, id games.GameID) (*model.Game, error) {
	// Mock implementation
	if id == "1" {
		return &model.Game{ID: "1"}, nil
	} else if id == "" {
		return nil, svcerrors.ErrInvalidID
	} else {
		return nil, svcerrors.ErrNotFound
	}
}

func (m *MockSingleController) Create(ctx context.Context, g *model.Game) (*model.Game, error) {
	// Mock implementation
	if g == nil {
		return nil, errors.New("the game to be created cannot be nil")
	}
	g.ID = "1" // Assign a mock ID
	return g, nil
}

func (m *MockSingleController) Replace(ctx context.Context, g *model.Game) (*model.Game, error) {
	// Mock implementation
	if g == nil {
		return nil, errors.New("the game to be replaced cannot be nil")
	}
	if g.ID == "" {
		return nil, svcerrors.ErrInvalidID
	}
	if g.ID != "1" {
		return nil, svcerrors.ErrNotFound
	}
	return g, nil
}

func (m *MockSingleController) DeleteByID(ctx context.Context, id games.GameID) (bool, error) {
	// Mock implementation
	if id == "" {
		return false, svcerrors.ErrInvalidID
	}
	if id == "1" {
		return true, nil
	}
	return false, svcerrors.ErrNotFound
}

func TestHumaHandlerGetByID(t *testing.T) {
	handler := createHumaHandler()

	// Test valid ID
	res, err := handler.GetByID(context.Background(), &GetByIDRequest{ID: "1"})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if res.Body.ID != "1" {
		t.Fatalf("expected game ID '1', got '%s'", res.Body.ID)
	}

	// Test invalid ID
	_, err = handler.GetByID(context.Background(), &GetByIDRequest{ID: "999"})
	if err == nil {
		t.Fatalf("expected error for non-existent game, got nil")
	} else if err.GetStatus() != 404 {
		t.Fatalf("expected 404 for non-existent game, got %v", err)
	}

	// Test empty ID
	_, err = handler.GetByID(context.Background(), &GetByIDRequest{ID: ""})
	if err == nil {
		t.Fatalf("expected error for empty game ID, got nil")
	} else if err.GetStatus() != 400 {
		t.Fatalf("expected 400 for empty game ID, got %v", err)
	}
}

func createHumaHandler() *HumaHandler {
	return NewHumaHandler(&MockSingleController{})
}
