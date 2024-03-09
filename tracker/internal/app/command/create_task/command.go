package create_task

import (
	"context"
	"fmt"
	"math/rand"
	"strings"

	uuid "github.com/satori/go.uuid"
	"github.com/segmentio/kafka-go"
	"google.golang.org/protobuf/proto"

	"github.com/ilyatos/aac-task-tracker/schema_registry/pkg/meta"
	task_event "github.com/ilyatos/aac-task-tracker/schema_registry/pkg/task"
	"github.com/ilyatos/aac-task-tracker/schema_registry/pkg/task/task_created"
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

func (c *Command) Handle(ctx context.Context, description string) (uuid.UUID, error) {
	users, err := c.userRepository.GetAllByRole(ctx, user_model.Employee)
	if err != nil {
		return uuid.UUID{}, fmt.Errorf("get users error: %w", err)
	}
	if users == nil || len(users) == 0 {
		return uuid.UUID{}, fmt.Errorf("users not found")
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
		return uuid.UUID{}, fmt.Errorf("create task error: %w", err)
	}

	// TODO: outbox pattern
	err = produceTaskCreatedEvent(ctx, newTask)
	if err != nil {
		return uuid.UUID{}, fmt.Errorf("produce task created event error: %w", err)
	}

	return newTask.PublicID, nil
}

func produceTaskCreatedEvent(ctx context.Context, task repository.NewTask) error {
	taskCreatedEvent, err := proto.Marshal(&task_created.TaskCreated{
		Header: &meta.Header{
			Producer: "tracker.create_task",
		},
		Payload: &task_created.TaskCreated_V1{
			V1: &task_created.V1{
				PublicId:     task.PublicID.String(),
				UserPublicId: task.UserPublicID.String(),
				Description:  task.Description,
				Status:       task_event.Status(task_event.Status_value[strings.ToUpper(string(task.Status))]),
			},
		},
	})
	if err != nil {
		return err
	}

	err = broker.Produce(ctx, "tasks-stream", kafka.Message{
		Key:   []byte("TaskCreated"),
		Value: taskCreatedEvent,
	})

	return err
}
