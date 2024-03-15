package get_most_expensive_task

import (
	"context"

	"github.com/ilyatos/aac-task-tracker/analytics/gen/analytics"
)

type Handler struct {
}

func New() *Handler {
	return &Handler{}
}

func (h *Handler) Handle(ctx context.Context, request *analytics.GetMostExpensiveTaskRequest) (*analytics.GetMostExpensiveTaskResponse, error) {
	return &analytics.GetMostExpensiveTaskResponse{}, nil
}
