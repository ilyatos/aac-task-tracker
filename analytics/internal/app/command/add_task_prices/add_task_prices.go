package add_task_prices

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

type AddTaskPricesData struct {
	TaskPublicID uuid.UUID
	CreditPrice  uint64
	DebitPrice   uint64
}

func (c *Command) Handle(ctx context.Context, addTaskPricesData AddTaskPricesData) error {
	task, err := c.taskRepository.GetByPublicID(ctx, addTaskPricesData.TaskPublicID)
	if err != nil {
		return fmt.Errorf("get task by public id error: %w", err)
	}

	task.CreditPrice = &addTaskPricesData.CreditPrice
	task.DebitPrice = &addTaskPricesData.DebitPrice

	err = c.taskRepository.Update(ctx, task)
	if err != nil {
		return fmt.Errorf("add task prices error: %w", err)
	}

	return nil
}
