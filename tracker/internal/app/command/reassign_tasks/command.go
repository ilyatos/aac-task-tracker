package reassign_tasks

import (
	"context"
	"encoding/json"
	"fmt"
	"math/rand"

	uuid "github.com/satori/go.uuid"
	"github.com/segmentio/kafka-go"

	task_model "github.com/ilyatos/aac-task-tracker/tracker/internal/app/model/task"
	user_model "github.com/ilyatos/aac-task-tracker/tracker/internal/app/model/user"
	"github.com/ilyatos/aac-task-tracker/tracker/internal/app/repository"
	"github.com/ilyatos/aac-task-tracker/tracker/internal/broker"
)

type Command struct {
	userRepository *repository.UserRepository
	taskRepository *repository.TaskRepository
}

func New(userRepository *repository.UserRepository, taskRepository *repository.TaskRepository) *Command {
	return &Command{userRepository: userRepository, taskRepository: taskRepository}
}

type taskAssignedEvent struct {
	PublicID     uuid.UUID `json:"public_id"`
	UserPublicID uuid.UUID `json:"user_public_id"`
}

type taskUpdatedEvent struct {
	PublicID     uuid.UUID         `json:"public_id"`
	UserPublicID uuid.UUID         `json:"user_public_id"`
	Description  string            `json:"description"`
	Status       task_model.Status `json:"status"`
}

func (c *Command) Handle(ctx context.Context) error {
	users, err := c.userRepository.GetAllByRole(ctx, user_model.Employee)
	if err != nil {
		return fmt.Errorf("get users error: %w", err)
	}
	if users == nil || len(users) == 0 {
		return fmt.Errorf("users not found")
	}

	tasks, err := c.taskRepository.GetAllByStatus(ctx, task_model.New)
	if err != nil {
		return fmt.Errorf("get tasks error: %w", err)
	}
	if tasks == nil || len(tasks) == 0 {
		return fmt.Errorf("tasks not found")
	}

	// TODO: parallelize with wait group
	for _, task := range tasks {
		task.UserPublicID = users[rand.Intn(len(users))].PublicID

		err = c.taskRepository.Update(ctx, task)
		if err != nil {
			return fmt.Errorf("update task error: %w", err)
		}

		err = produceTaskAssignedEvent(ctx, taskAssignedEvent{
			PublicID:     task.PublicID,
			UserPublicID: task.UserPublicID,
		})
		if err != nil {
			return fmt.Errorf("produce task created event error: %w", err)
		}

		err = produceTaskUpdatedEvent(ctx, taskUpdatedEvent(task))
		if err != nil {
			return fmt.Errorf("produce task updated event error: %w", err)
		}
	}

	return nil
}

func produceTaskAssignedEvent(ctx context.Context, taskAssignedEvent taskAssignedEvent) error {
	taskAssignedValue, err := json.Marshal(taskAssignedEvent)
	if err != nil {
		return err
	}

	err = broker.Produce(ctx, "tasks", kafka.Message{
		Key:   []byte("TaskAssigned"),
		Value: taskAssignedValue,
	})

	return err
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
