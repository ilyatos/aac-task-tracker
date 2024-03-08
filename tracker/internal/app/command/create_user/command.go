package create_user

import (
	"context"
	"fmt"

	uuid "github.com/satori/go.uuid"

	user_model "github.com/ilyatos/aac-task-tracker/tracker/internal/app/model/user"
	"github.com/ilyatos/aac-task-tracker/tracker/internal/app/repository"
)

type Command struct {
	userRepository *repository.UserRepository
}

func New(userRepository *repository.UserRepository) *Command {
	return &Command{userRepository: userRepository}
}

type CreateUserData struct {
	PublicID uuid.UUID
	Name     string
	Email    string
	Role     user_model.Role
}

func (c *Command) Handle(ctx context.Context, createUserData CreateUserData) error {
	err := c.userRepository.Create(ctx, repository.User(createUserData))
	if err != nil {
		return fmt.Errorf("create user error: %w", err)
	}

	return nil
}
