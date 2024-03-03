package change_role

import (
	"context"
	"fmt"

	uuid "github.com/satori/go.uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/ilyatos/aac-task-tracker/auth/gen/auth"
	"github.com/ilyatos/aac-task-tracker/auth/internal/app/command/change_role"
	"github.com/ilyatos/aac-task-tracker/auth/internal/app/model/user"
)

var rpcRoleToUserRoleMap = map[auth.Role]user.Role{
	auth.Role_EMPLOYEE:   user.Employee,
	auth.Role_MANAGER:    user.Manager,
	auth.Role_ACCOUNTANT: user.Accountant,
	auth.Role_ADMIN:      user.Admin,
}

type Handler struct {
	changeRole *change_role.Command
}

func New(changeRole *change_role.Command) *Handler {
	return &Handler{changeRole: changeRole}
}

func (h *Handler) Handle(ctx context.Context, request *auth.ChangeRoleRequest) (*auth.ChangeRoleResponse, error) {
	userId, err := uuid.FromString(request.GetUserId())
	if err != nil {
		return &auth.ChangeRoleResponse{}, status.Error(codes.Internal, fmt.Sprintf("invalid user id: %s", err))
	}

	err = h.changeRole.Handle(ctx, userId, rpcRoleToUserRoleMap[request.GetRole()])
	if err != nil {
		return &auth.ChangeRoleResponse{}, status.Error(codes.Internal, fmt.Sprintf("changing role failed: %s", err))
	}

	return &auth.ChangeRoleResponse{}, nil
}
