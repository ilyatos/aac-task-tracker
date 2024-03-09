package complete_task

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	uuid "github.com/satori/go.uuid"
	"github.com/segmentio/kafka-go"
	"google.golang.org/protobuf/proto"

	"github.com/ilyatos/aac-task-tracker/schema_registry/pkg/meta"
	task_event "github.com/ilyatos/aac-task-tracker/schema_registry/pkg/task"
	"github.com/ilyatos/aac-task-tracker/schema_registry/pkg/task/task_completed"
	"github.com/ilyatos/aac-task-tracker/schema_registry/pkg/task/task_updated"
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

	err = produceTaskCompletedEvent(ctx, task)
	if err != nil {
		return fmt.Errorf("produce task completed event error: %w", err)
	}

	err = produceTaskUpdatedEvent(ctx, task)
	if err != nil {
		return fmt.Errorf("produce task updated event error: %w", err)
	}

	return nil
}

func produceTaskCompletedEvent(ctx context.Context, task *repository.Task) error {
	taskCompleted, err := proto.Marshal(&task_completed.TaskCompleted{
		Header:  &meta.Header{Producer: "tracker.complete_task"},
		Payload: &task_completed.TaskCompleted_V1{V1: &task_completed.V1{PublicId: task.PublicID.String()}},
	})
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

func produceTaskUpdatedEvent(ctx context.Context, task *repository.Task) error {
	taskUpdated, err := json.Marshal(&task_updated.TaskUpdated{
		Header: &meta.Header{Producer: "tracker.complete_task"},
		Payload: &task_updated.TaskUpdated_V1{
			V1: &task_updated.V1{
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
		Key:   []byte("TaskUpdated"),
		Value: taskUpdated,
	})
	if err != nil {
		return err
	}

	return nil
}
