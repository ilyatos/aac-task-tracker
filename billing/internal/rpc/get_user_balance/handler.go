package get_user_balance

import (
	"context"
	"fmt"

	"github.com/ilyatos/aac-task-tracker/billing/gen/billing"

	"github.com/ilyatos/aac-task-tracker/billing/internal/app/query/get_user_balance"
)

type Handler struct {
	getUserBalance *get_user_balance.Query
}

func New(getUserBalance *get_user_balance.Query) *Handler {
	return &Handler{getUserBalance: getUserBalance}
}

func (h *Handler) Handle(ctx context.Context, request *billing.GetUserBalanceRequest) (*billing.GetUserBalanceResponse, error) {
	err := h.getUserBalance.Handle(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get user balance: %w", err)
	}

	return &billing.GetUserBalanceResponse{}, nil
}
