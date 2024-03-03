package create_task

import (
	"context"
	"fmt"

	"github.com/ilyatos/aac-task-tracker/tracker/gen/tracker"
	"github.com/ilyatos/aac-task-tracker/tracker/internal/app/command/create_task"
)

type Handler struct {
	createTask *create_task.Command
}

func New(createTask *create_task.Command) *Handler {
	return &Handler{createTask: createTask}
}

func (h *Handler) Handle(ctx context.Context, request *tracker.CreateTaskRequest) (*tracker.CreateTaskResponse, error) {
	err := h.createTask.Handle(ctx, request.GetDescription())
	if err != nil {
		return nil, fmt.Errorf("failed to create task: %w", err)
	}

	return &tracker.CreateTaskResponse{}, nil
}
