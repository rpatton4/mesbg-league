package client

import (
	"context"
	"encoding/json"
	"fmt"
	gamesheader "github.com/rpatton4/mesbg-league/games/pkg/header"
	games "github.com/rpatton4/mesbg-league/games/pkg/model"
	"github.com/rpatton4/mesbg-league/pkg/svcerrors"
	"net/http"
)

type HTTPGateway struct {
	addr string
}

func New(addr string) *HTTPGateway {
	return &HTTPGateway{addr: addr}
}

func (g *HTTPGateway) GetByID(ctx context.Context, id gamesheader.GameID) (*games.Game, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, g.addr+"/"+string(id), nil)
	if err != nil {
		return nil, err
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNotFound {
		return nil, svcerrors.NotFound
	} else if resp.StatusCode/100 != 2 {
		return nil, fmt.Errorf("non-2xx response: %v", resp)
	}
	var game *games.Game
	if err := json.NewDecoder(resp.Body).Decode(&game); err != nil {
		return nil, fmt.Errorf("failed to decode game: %w", err)
	}
	return game, nil
}
