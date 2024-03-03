package login

import (
	"context"
	"fmt"

	"github.com/ilyatos/aac-task-tracker/auth/internal/app/helper/jwt"
	"github.com/ilyatos/aac-task-tracker/auth/internal/app/helper/password"
	"github.com/ilyatos/aac-task-tracker/auth/internal/app/repository"
)

type Command struct {
	userRepository *repository.UserRepository
	hasher         *password.Hasher
	jwt            *jwt.Manager
}

func New(userRepository *repository.UserRepository, hasher *password.Hasher, jwt *jwt.Manager) *Command {
	return &Command{userRepository: userRepository, hasher: hasher, jwt: jwt}
}

func (l *Command) Handle(ctx context.Context, email, password string) (string, error) {
	user, err := l.userRepository.FindByEmail(ctx, email)
	if err != nil {
		return "", fmt.Errorf("find by email error: %w", err)
	}
	if user == nil {
		return "", fmt.Errorf("invalid user credentials")
	}

	if r := l.hasher.CheckHash(password, user.Password); !r {
		return "", fmt.Errorf("invalid user credentials")
	}

	return l.jwt.NewAccessToken(jwt.NewUserClaims(user.PublicID))
}
