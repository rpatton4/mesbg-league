// The Games gateway package contains clients for interacting with the gmes service from other services.
// The package is meant to be public, and will commmit to be backwards compatible within major versions.
// Te package primarily consists of the GamesGateway interface, with sub-packages holding different implementations
// of the interface for calling it in-memory or over HTTP.
package client

import (
	"context"
	gamesheader "github.com/rpatton4/mesbg-league/games/pkg/header"
	games "github.com/rpatton4/mesbg-league/games/pkg/model"
)

// GamesGateway provides a set of methods for interacting with the Games service from outside the service.
type GamesGateway interface {
	GetByID(ctx context.Context, id gamesheader.GameID) (*games.Game, error)
	Create(ctx context.Context, game *games.Game) (*games.Game, error)
}
