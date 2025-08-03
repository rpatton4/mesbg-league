package gateway

import (
	"context"
	"github.com/rpatton4/mesbg-league/games/internal/primary"
	"github.com/rpatton4/mesbg-league/games/internal/secondary"
	games "github.com/rpatton4/mesbg-league/games/pkg"
	"github.com/rpatton4/mesbg-league/games/pkg/model"
)

type InProcessGateway struct {
	ctrl primary.SingleController
}

// NewInProcessGatewayWithController creates a new InProcessGateway with the provided transaction controller.
// This is intended primarily for use while testing, to provide a mock or stub controller.
func NewInProcessGatewayWithController(ctrl primary.SingleController) *InProcessGateway {
	return &InProcessGateway{ctrl}
}

// NewDefaultInProcessGateway creates a new InProcessGateway with a default repository and txncontroller
func NewDefaultInProcessGateway() *InProcessGateway {
	repo := secondary.NewDefaultRepository()
	ctrl := primary.NewTxnController(repo)
	return NewInProcessGatewayWithController(ctrl)
}

func (ipg *InProcessGateway) GetByID(ctx context.Context, id games.GameID) (*model.Game, error) {
	return ipg.ctrl.GetByID(ctx, id)
}
func (ipg *InProcessGateway) Create(ctx context.Context, g *model.Game) (*model.Game, error) {
	return ipg.ctrl.Create(ctx, g)
}
func (ipg *InProcessGateway) Replace(ctx context.Context, g *model.Game) (*model.Game, error) {
	return ipg.ctrl.Replace(ctx, g)
}
func (ipg *InProcessGateway) DeleteByID(ctx context.Context, id games.GameID) (bool, error) {
	return ipg.ctrl.DeleteByID(ctx, id)
}
