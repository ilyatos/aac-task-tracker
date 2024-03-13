package create_task

import (
	"context"
	"fmt"
	"math/rand"

	uuid "github.com/satori/go.uuid"

	"github.com/ilyatos/aac-task-tracker/billing/internal/app/repository"
)

type Command struct {
	taskRepository *repository.TaskRepository
}

func New(taskRepository *repository.TaskRepository) *Command {
	return &Command{taskRepository: taskRepository}
}

type CreateTaskData struct {
	PublicID     uuid.UUID
	UserPublicID uuid.UUID
	Description  string
}

func (c *Command) Handle(ctx context.Context, createTaskData CreateTaskData) error {
	newTask := repository.Task{
		PublicID:     createTaskData.PublicID,
		UserPublicID: createTaskData.UserPublicID,
		Description:  createTaskData.Description,
		CreditPrice:  getCreditPrice(),
		DebitPrice:   getDebitPrice(),
	}

	err := c.taskRepository.Create(ctx, newTask)
	if err != nil {
		return fmt.Errorf("create task error: %w", err)
	}

	return nil
}

func getCreditPrice() uint64 {
	return uint64(rand.Intn(40-20) + 20)
}

func getDebitPrice() uint64 {
	return uint64(rand.Intn(20-10) + 10)
}
