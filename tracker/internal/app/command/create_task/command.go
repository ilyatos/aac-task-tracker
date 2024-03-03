package create_task

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

type taskCreatedEvent struct {
	PublicID     uuid.UUID         `json:"public_id"`
	UserPublicID uuid.UUID         `json:"user_public_id"`
	Description  string            `json:"description"`
	Status       task_model.Status `json:"status"`
}

func (c *Command) Handle(ctx context.Context, description string) error {
	users, err := c.userRepository.GetAllByRole(ctx, user_model.Employee)
	if err != nil {
		return fmt.Errorf("get users error: %w", err)
	}
	if users == nil || len(users) == 0 {
		return fmt.Errorf("users not found")
	}

	user := users[rand.Intn(len(users))]

	newTask := repository.NewTask{
		PublicID:     uuid.NewV4(),
		UserPublicID: user.PublicID,
		Description:  description,
		Status:       task_model.New,
	}

	err = c.taskRepository.Create(ctx, newTask)
	if err != nil {
		return fmt.Errorf("create task error: %w", err)
	}

	err = produceTaskCreatedEvent(ctx, taskCreatedEvent(newTask))
	if err != nil {
		return fmt.Errorf("produce task created event error: %w", err)
	}

	return nil
}

func produceTaskCreatedEvent(ctx context.Context, taskCreatedEvent taskCreatedEvent) error {
	userCreatedEvent, err := json.Marshal(taskCreatedEvent)
	if err != nil {
		return err
	}

	err = broker.Produce(ctx, "tasks-stream", kafka.Message{
		Key:   []byte("TaskCreated"),
		Value: userCreatedEvent,
	})

	return err
}
