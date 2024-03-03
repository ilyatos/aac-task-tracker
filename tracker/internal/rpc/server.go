package rpc

import (
	"context"

	"github.com/ilyatos/aac-task-tracker/tracker/gen/tracker"
)

type (
	createTask    func(ctx context.Context, request *tracker.CreateTaskRequest) (*tracker.CreateTaskResponse, error)
	reassignTasks func(ctx context.Context, request *tracker.ReassignTasksRequest) (*tracker.ReassignTasksResponse, error)
	completeTask  func(ctx context.Context, request *tracker.CompleteTaskRequest) (*tracker.CompleteTaskResponse, error)
)

// Server unions standalone handlers.
type Server struct {
	tracker.TrackerServer
	createTask    createTask
	reassignTasks reassignTasks
	completeTask  completeTask
}

// New returns Server instance.
func New(
	createTask createTask,
	reassignTasks reassignTasks,
	completeTask completeTask,
) *Server {
	return &Server{
		createTask:    createTask,
		reassignTasks: reassignTasks,
		completeTask:  completeTask,
	}
}

func (s Server) CreateTask(ctx context.Context, request *tracker.CreateTaskRequest) (*tracker.CreateTaskResponse, error) {
	return s.createTask(ctx, request)
}

func (s Server) ReassignTasks(ctx context.Context, request *tracker.ReassignTasksRequest) (*tracker.ReassignTasksResponse, error) {
	return s.reassignTasks(ctx, request)
}

func (s Server) CompleteTask(ctx context.Context, request *tracker.CompleteTaskRequest) (*tracker.CompleteTaskResponse, error) {
	return s.completeTask(ctx, request)
}
