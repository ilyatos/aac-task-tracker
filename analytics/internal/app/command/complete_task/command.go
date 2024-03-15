package complete_task

import (
	"context"
	"time"

	uuid "github.com/satori/go.uuid"

	"github.com/ilyatos/aac-task-tracker/analytics/internal/app/repository"
)

type Command struct {
	taskRepository *repository.TaskRepository
}

func New(taskRepository *repository.TaskRepository) *Command {
	return &Command{taskRepository: taskRepository}
}

type CompleteTaskData struct {
	PublicID    uuid.UUID
	CompletedAt time.Time
}

func (c *Command) Handle(ctx context.Context, completeTaskData CompleteTaskData) error {
	task, err := c.taskRepository.GetByPublicID(ctx, completeTaskData.PublicID)
	if err != nil {
		return err
	}

	task.CompletedAt = &completeTaskData.CompletedAt

	err = c.taskRepository.Update(ctx, task)
	if err != nil {
		return err
	}

	return nil
}
