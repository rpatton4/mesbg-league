package gateway

import (
	"context"
	gamesctrl "github.com/rpatton4/mesbg-league/games/internal/gamesprimaryports"
	gamesheader "github.com/rpatton4/mesbg-league/games/pkg/header"
	gamesmodel "github.com/rpatton4/mesbg-league/games/pkg/model"
)

type InProcessGateway struct {
	ctrl *gamesctrl.TxnController
}

func NewInProcessGateway(ctrl *gamesctrl.TxnController) *InProcessGateway {
	return &InProcessGateway{ctrl}
}

func (ipg *InProcessGateway) GetByID(ctx context.Context, id gamesheader.GameID) (*gamesmodel.Game, error) {
	return ipg.ctrl.GetByID(ctx, id)
}
func (ipg *InProcessGateway) Create(ctx context.Context, g *gamesmodel.Game) (*gamesmodel.Game, error) {
	return ipg.ctrl.Create(ctx, g)
}
