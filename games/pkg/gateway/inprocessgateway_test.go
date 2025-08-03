package gateway

import (
	mock_primary "github.com/rpatton4/mesbg-league/games/internal/primary/mocks"
	"go.uber.org/mock/gomock"
	"testing"
)

// TestInProcessGatewayMockedGetByID tests the GetByID method of the InProcessGateway using a mocked controller.
// In most respects this mirrors the testing done in the humahandler_test.go file for the HumaHandler's GetByID method,
func TestInProcessGatewayMockedGetByID(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	mockController := mock_primary.NewMockSingleController(mockCtrl)
	gtwy := NewInProcessGatewayWithController(mockController)

}
