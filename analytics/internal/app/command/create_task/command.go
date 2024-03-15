package create_user

import (
	"context"
	"fmt"

	uuid "github.com/satori/go.uuid"

	"github.com/ilyatos/aac-task-tracker/analytics/internal/app/repository"
)

type Command struct {
	taskRepository *repository.TaskRepository
}

func New(taskRepository *repository.TaskRepository) *Command {
	return &Command{taskRepository: taskRepository}
}

type CreateTaskData struct {
	PublicID    uuid.UUID
	Description string
}

func (c *Command) Handle(ctx context.Context, createTaskData CreateTaskData) error {
	err := c.taskRepository.Create(ctx, repository.Task{
		PublicID:    createTaskData.PublicID,
		Description: createTaskData.Description,
		CreditPrice: nil,
		DebitPrice:  nil,
		CompletedAt: nil,
	})
	if err != nil {
		return fmt.Errorf("create user error: %w", err)
	}

	return nil
}
