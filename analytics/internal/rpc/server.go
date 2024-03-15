package rpc

import (
	"context"

	"github.com/ilyatos/aac-task-tracker/analytics/gen/analytics"
)

type (
	getMostExpensiveTask func(ctx context.Context, request *analytics.GetMostExpensiveTaskRequest) (*analytics.GetMostExpensiveTaskResponse, error)
)

// Server unions standalone handlers.
type Server struct {
	analytics.AnalyticsServer
	getMostExpensiveTask getMostExpensiveTask
}

// New returns Server instance.
func New(
	getMostExpensiveTask getMostExpensiveTask,
) *Server {
	return &Server{
		getMostExpensiveTask: getMostExpensiveTask,
	}
}

func (s Server) GetMostExpensiveTask(ctx context.Context, request *analytics.GetMostExpensiveTaskRequest) (*analytics.GetMostExpensiveTaskResponse, error) {
	return s.getMostExpensiveTask(ctx, request)
}
