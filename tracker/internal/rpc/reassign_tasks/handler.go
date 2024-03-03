package reassign_tasks

import (
	"context"
	"fmt"

	"github.com/ilyatos/aac-task-tracker/tracker/gen/tracker"
	"github.com/ilyatos/aac-task-tracker/tracker/internal/app/command/reassign_tasks"
)

type Handler struct {
	reassignTasks *reassign_tasks.Command
}

func New(reassignTasks *reassign_tasks.Command) *Handler {
	return &Handler{reassignTasks: reassignTasks}
}

func (h *Handler) Handle(ctx context.Context, request *tracker.ReassignTasksRequest) (*tracker.ReassignTasksResponse, error) {
	err := h.reassignTasks.Handle(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to reassign tasks: %w", err)
	}

	return &tracker.ReassignTasksResponse{}, nil
}
