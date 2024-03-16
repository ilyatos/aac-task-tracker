package tasks_stream

import (
	"context"
	"fmt"
	"log"
	"os"

	uuid "github.com/satori/go.uuid"
	"github.com/segmentio/kafka-go"
	"google.golang.org/protobuf/proto"

	create_task_command "github.com/ilyatos/aac-task-tracker/billing/internal/app/command/create_task"
	update_task_command "github.com/ilyatos/aac-task-tracker/billing/internal/app/command/update_task"
	"github.com/ilyatos/aac-task-tracker/schema_registry/pkg/task/task_created"
	"github.com/ilyatos/aac-task-tracker/schema_registry/pkg/task/task_updated"
)

type Consumer struct {
	reader            *kafka.Reader
	createTaskCommand *create_task_command.Command
	updateTaskCommand *update_task_command.Command
}

func New(
	createTaskCommand *create_task_command.Command,
	updateTaskCommand *update_task_command.Command,
) *Consumer {
	r := kafka.NewReader(kafka.ReaderConfig{
		Brokers:  []string{os.Getenv("KAFKA_URL")},
		GroupID:  "billing-group",
		Topic:    "tasks-stream",
		MaxBytes: 10e6, // 10MB
	})

	return &Consumer{reader: r, createTaskCommand: createTaskCommand, updateTaskCommand: updateTaskCommand}
}

func (c *Consumer) Consume() {
	log.Println("start consuming tasks stream")

	ctx := context.Background()
	for {
		m, err := c.reader.FetchMessage(ctx)
		if err != nil {
			log.Println("failed to commit messages:", err)
			continue
		}

		log.Printf("message at topic/partition/offset %v/%v/%v: %s = %s\n", m.Topic, m.Partition, m.Offset, string(m.Key), string(m.Value))

		switch string(m.Key) {
		case "TaskCreated":
			err = c.handleTaskCreated(ctx, m.Value)
		case "TaskUpdated":
			err = c.handleTaskUpdated(ctx, m.Value)
		default:
			log.Println("unknown message key:", string(m.Key))
		}

		if err != nil {
			log.Println("failed to handle message:", err)
			continue
		}

		if err := c.reader.CommitMessages(ctx, m); err != nil {
			log.Println("failed to commit messages:", err)
			continue
		}
	}
}

func (c *Consumer) handleTaskCreated(ctx context.Context, event []byte) error {
	taskCreated := task_created.TaskCreated{}
	err := proto.Unmarshal(event, &taskCreated)
	if err != nil {
		return fmt.Errorf("unmarshal message error: %w", err)
	}

	switch taskCreated.Payload.(type) {
	case *task_created.TaskCreated_V1:
		v1 := taskCreated.GetV1()
		return c.createTaskCommand.Handle(ctx, create_task_command.CreateTaskData{
			PublicID:     uuid.FromStringOrNil(v1.GetPublicId()),
			UserPublicID: uuid.FromStringOrNil(v1.GetUserPublicId()),
			Description:  v1.GetDescription(),
		})
	}

	return nil
}

func (c *Consumer) handleTaskUpdated(ctx context.Context, event []byte) error {
	taskUpdate := task_updated.TaskUpdated{}
	err := proto.Unmarshal(event, &taskUpdate)
	if err != nil {
		return fmt.Errorf("unmarshal message error: %w", err)
	}

	switch taskUpdate.Payload.(type) {
	case *task_updated.TaskUpdated_V1:
		v1 := taskUpdate.GetV1()
		return c.updateTaskCommand.Handle(ctx, update_task_command.UpdateTaskData{
			PublicID:     uuid.FromStringOrNil(v1.GetPublicId()),
			UserPublicID: uuid.FromStringOrNil(v1.GetUserPublicId()),
			Description:  v1.GetDescription(),
		})
	}

	return nil
}
