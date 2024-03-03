package signup

import (
	"context"
	"fmt"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/ilyatos/aac-task-tracker/auth/gen/auth"
	"github.com/ilyatos/aac-task-tracker/auth/internal/app/command/singup"
	"github.com/ilyatos/aac-task-tracker/auth/internal/app/model/user"
)

var rpcRoleToUserRoleMap = map[auth.Role]user.Role{
	auth.Role_EMPLOYEE:   user.Employee,
	auth.Role_MANAGER:    user.Manager,
	auth.Role_ACCOUNTANT: user.Accountant,
	auth.Role_ADMIN:      user.Admin,
}

type Handler struct {
	signup *singup.Command
}

func New(signup *singup.Command) *Handler {
	return &Handler{signup: signup}
}

func (h *Handler) Handle(ctx context.Context, request *auth.SignUpRequest) (*auth.SignUpResponse, error) {
	err := h.signup.Handle(ctx, singup.UserSignUpData{
		Name:     request.GetName(),
		Email:    request.GetEmail(),
		Password: request.GetPassword(),
		Role:     rpcRoleToUserRoleMap[request.GetRole()],
	})
	if err != nil {
		fmt.Println("signup failed:", err)
		return &auth.SignUpResponse{}, status.Error(codes.Internal, fmt.Sprintf("signup failed: %s", err))
	}

	return &auth.SignUpResponse{}, nil
}
