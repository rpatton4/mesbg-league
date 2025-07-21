package client

import (
	"context"
	gamesctrl "github.com/rpatton4/mesbg-league/games/internal/inbound"
	gamesheader "github.com/rpatton4/mesbg-league/games/pkg/header"
	games "github.com/rpatton4/mesbg-league/games/pkg/model"
)

type InProcessGateway struct {
	ctrl *gamesctrl.Controller
}

func New(ctrl *gamesctrl.Controller) *InProcessGateway {
	return &InProcessGateway{ctrl}
}

func (g *InProcessGateway) GetByID(ctx context.Context, id gamesheader.GameID) (*games.Game, error) {
	return g.ctrl.GetByID(ctx, id)
}
