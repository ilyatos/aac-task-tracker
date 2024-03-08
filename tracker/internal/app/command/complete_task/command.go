package complete_task

import (
	"context"
	"encoding/json"
	"fmt"

	uuid "github.com/satori/go.uuid"
	"github.com/segmentio/kafka-go"

	task_model "github.com/ilyatos/aac-task-tracker/tracker/internal/app/model/task"
	"github.com/ilyatos/aac-task-tracker/tracker/internal/app/repository"
	"github.com/ilyatos/aac-task-tracker/tracker/internal/broker"
)

type Command struct {
	taskRepository *repository.TaskRepository
}

func New(taskRepository *repository.TaskRepository) *Command {
	return &Command{taskRepository: taskRepository}
}

type taskCompletedEvent struct {
	PublicID uuid.UUID `json:"public_id"`
}

type taskUpdatedEvent struct {
	PublicID     uuid.UUID         `json:"public_id"`
	UserPublicID uuid.UUID         `json:"user_public_id"`
	Description  string            `json:"description"`
	Status       task_model.Status `json:"status"`
}

func (c *Command) Handle(ctx context.Context, userPublicID uuid.UUID, taskPublicID uuid.UUID) error {
	task, err := c.taskRepository.GetByPublicID(ctx, taskPublicID)
	if err != nil {
		return fmt.Errorf("get task error: %w", err)
	}
	if task == nil {
		return fmt.Errorf("task not found")
	}
	if task.UserPublicID != userPublicID {
		return fmt.Errorf("forbidden: user is not the owner of the task")
	}

	task.Status = task_model.Done

	err = c.taskRepository.Update(ctx, *task)
	if err != nil {
		return fmt.Errorf("update task error: %w", err)
	}

	err = produceTaskCompletedEvent(ctx, taskCompletedEvent{PublicID: taskPublicID})
	if err != nil {
		return fmt.Errorf("produce task completed event error: %w", err)
	}

	err = produceTaskUpdatedEvent(ctx, taskUpdatedEvent(*task))
	if err != nil {
		return fmt.Errorf("produce task updated event error: %w", err)
	}

	return nil
}

func produceTaskCompletedEvent(ctx context.Context, taskCompletedEvent taskCompletedEvent) error {
	taskCompleted, err := json.Marshal(taskCompletedEvent)
	if err != nil {
		return err
	}

	err = broker.Produce(ctx, "tasks", kafka.Message{
		Key:   []byte("TaskCompleted"),
		Value: taskCompleted,
	})
	if err != nil {
		return err
	}

	return nil
}

func produceTaskUpdatedEvent(ctx context.Context, taskUpdatedEvent taskUpdatedEvent) error {
	taskUpdated, err := json.Marshal(taskUpdatedEvent)
	if err != nil {
		return err
	}

	err = broker.Produce(ctx, "tasks-stream", kafka.Message{
		Key:   []byte("TaskUpdated"),
		Value: taskUpdated,
	})
	if err != nil {
		return err
	}

	return nil
}
