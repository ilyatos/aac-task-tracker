package reassign_tasks

import (
	"context"
	"encoding/json"

	uuid "github.com/satori/go.uuid"
	"github.com/segmentio/kafka-go"

	task_model "github.com/ilyatos/aac-task-tracker/tracker/internal/app/model/task"
	"github.com/ilyatos/aac-task-tracker/tracker/internal/broker"
)

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

func makeTaskAssignedMessage(ctx context.Context, taskAssignedEvent taskAssignedEvent) (kafka.Message, error) {
	taskAssignedValue, err := json.Marshal(taskAssignedEvent)
	if err != nil {
		return kafka.Message{}, err
	}

	return kafka.Message{
		Key:   []byte("TaskAssigned"),
		Value: taskAssignedValue,
	}, nil
}

func produceTaskAssignedEvents(ctx context.Context, msgs []kafka.Message) error {
	err := broker.Produce(ctx, "tasks", msgs...)
	if err != nil {
		return err
	}

	return nil
}

func makeTaskUpdatedMessage(ctx context.Context, taskUpdatedEvent taskUpdatedEvent) (kafka.Message, error) {
	taskUpdatedValue, err := json.Marshal(taskUpdatedEvent)
	if err != nil {
		return kafka.Message{}, err
	}

	return kafka.Message{
		Key:   []byte("TaskUpdated"),
		Value: taskUpdatedValue,
	}, nil
}

func produceTaskUpdatedEvents(ctx context.Context, msgs []kafka.Message) error {
	err := broker.Produce(ctx, "tasks-stream", msgs...)
	if err != nil {
		return err
	}

	return nil
}
