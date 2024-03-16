package create_user

import (
	"context"
	"fmt"

	"github.com/ilyatos/aac-task-tracker/analytics/internal/app/repository"
)

type Command struct {
	userRepository *repository.UserRepository
}

func New(userRepository *repository.UserRepository) *Command {
	return &Command{userRepository: userRepository}
}

type CreateUserData repository.User

func (c *Command) Handle(ctx context.Context, createUserData CreateUserData) error {
	err := c.userRepository.Create(ctx, repository.User{
		PublicID: createUserData.PublicID,
		Name:     createUserData.Name,
	})
	if err != nil {
		return fmt.Errorf("create user error: %w", err)
	}

	return nil
}
