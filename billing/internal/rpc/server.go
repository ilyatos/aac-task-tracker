package rpc

import (
	"context"

	"github.com/ilyatos/aac-task-tracker/billing/gen/billing"
)

type (
	getUserBalance func(ctx context.Context, request *billing.GetUserBalanceRequest) (*billing.GetUserBalanceResponse, error)
)

// Server unions standalone handlers.
type Server struct {
	billing.BillingServer
	getUserBalance getUserBalance
}

// New returns Server instance.
func New(
	getUserBalance getUserBalance,
) *Server {
	return &Server{
		getUserBalance: getUserBalance,
	}
}

func (s Server) GetUserBalance(ctx context.Context, request *billing.GetUserBalanceRequest) (*billing.GetUserBalanceResponse, error) {
	return s.getUserBalance(ctx, request)
}
