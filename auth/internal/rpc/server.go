package rpc

import (
	"context"
	"github.com/ilyatos/aac-task-tracker/auth/gen/auth"
)

type (
	singUp     func(ctx context.Context, request *auth.SignUpRequest) (*auth.SignUpResponse, error)
	logIn      func(ctx context.Context, request *auth.LogInRequest) (*auth.LogInResponse, error)
	changeRole func(ctx context.Context, request *auth.ChangeRoleRequest) (*auth.ChangeRoleResponse, error)
)

// Server unions standalone handlers.
type Server struct {
	auth.AuthServer
	singUp
	logIn
	changeRole
}

// New returns Server instance.
func New(
	singUp singUp,
	logIn logIn,
	changeRole changeRole,
) *Server {
	return &Server{
		singUp:     singUp,
		logIn:      logIn,
		changeRole: changeRole,
	}
}

func (s Server) SignUp(ctx context.Context, request *auth.SignUpRequest) (*auth.SignUpResponse, error) {
	return s.singUp(ctx, request)
}

func (s Server) LogIn(ctx context.Context, request *auth.LogInRequest) (*auth.LogInResponse, error) {
	return s.logIn(ctx, request)
}

func (s Server) ChangeRole(ctx context.Context, request *auth.ChangeRoleRequest) (*auth.ChangeRoleResponse, error) {
	return s.changeRole(ctx, request)
}
