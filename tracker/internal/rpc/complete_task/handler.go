package complete_task

import (
	"context"
	"fmt"

	"github.com/ilyatos/aac-task-tracker/tracker/internal/app/helper/jwt"
	"github.com/ilyatos/aac-task-tracker/tracker/internal/context_key"

	uuid "github.com/satori/go.uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/ilyatos/aac-task-tracker/tracker/gen/tracker"
	"github.com/ilyatos/aac-task-tracker/tracker/internal/app/command/complete_task"
)

type Handler struct {
	completeTask *complete_task.Command
}

func New(completeTask *complete_task.Command) *Handler {
	return &Handler{completeTask: completeTask}
}

func (h *Handler) Handle(ctx context.Context, request *tracker.CompleteTaskRequest) (*tracker.CompleteTaskResponse, error) {
	userClaims, ok := ctx.Value(context_key.AuthKey).(*jwt.UserClaims)
	if !ok || userClaims == nil {
		return &tracker.CompleteTaskResponse{}, status.Error(codes.Unauthenticated, "user claims not found")
	}

	taskPublicID, err := uuid.FromString(request.GetPublicId())
	if err != nil {
		return &tracker.CompleteTaskResponse{}, status.Error(codes.Internal, fmt.Sprintf("invalid task id: %s", err))
	}

	err = h.completeTask.Handle(ctx, userClaims.PublicID, taskPublicID)
	if err != nil {
		return nil, fmt.Errorf("failed to mark task as done: %w", err)
	}

	return &tracker.CompleteTaskResponse{}, nil
}
