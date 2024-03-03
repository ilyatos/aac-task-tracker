package login

import (
	"context"
	"fmt"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/ilyatos/aac-task-tracker/auth/gen/auth"
	"github.com/ilyatos/aac-task-tracker/auth/internal/app/command/login"
)

type Handler struct {
	login *login.Command
}

func New(login *login.Command) *Handler {
	return &Handler{login: login}
}

func (h *Handler) Handle(ctx context.Context, request *auth.LogInRequest) (*auth.LogInResponse, error) {
	jwt, err := h.login.Handle(ctx, request.Email, request.Password)
	if err != nil {
		fmt.Println("login failed:", err)
		return nil, status.Error(codes.Unauthenticated, fmt.Sprintf("login failed: %s", err))
	}

	return &auth.LogInResponse{Token: jwt}, nil
}
