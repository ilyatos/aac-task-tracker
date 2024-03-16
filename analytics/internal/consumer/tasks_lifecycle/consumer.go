package tasks_lifecycle

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	uuid "github.com/satori/go.uuid"
	"github.com/segmentio/kafka-go"
	"google.golang.org/protobuf/proto"

	complete_task_command "github.com/ilyatos/aac-task-tracker/analytics/internal/app/command/complete_task"
	"github.com/ilyatos/aac-task-tracker/schema_registry/pkg/task/task_completed"
)

type Consumer struct {
	reader              *kafka.Reader
	completeTaskCommand *complete_task_command.Command
}

func New(
	completeTaskCommand *complete_task_command.Command,
) *Consumer {
	r := kafka.NewReader(kafka.ReaderConfig{
		Brokers:  []string{os.Getenv("KAFKA_URL")},
		GroupID:  "analytics-group",
		Topic:    "tasks",
		MaxBytes: 10e6, // 10MB
	})

	return &Consumer{
		reader:              r,
		completeTaskCommand: completeTaskCommand,
	}
}

func (c *Consumer) Consume() {
	log.Println("start consuming tasks")

	ctx := context.Background()
	for {
		m, err := c.reader.FetchMessage(ctx)
		if err != nil {
			log.Println("failed to commit messages:", err)
			continue
		}

		log.Printf("message at topic/partition/offset %v/%v/%v: %s = %s\n", m.Topic, m.Partition, m.Offset, string(m.Key), string(m.Value))

		switch string(m.Key) {
		case "TaskCompleted":
			err = c.handleTaskCompleted(ctx, m.Value)
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

func (c *Consumer) handleTaskCompleted(ctx context.Context, event []byte) error {
	log.Println("handle task completed")

	taskCompleted := task_completed.TaskCompleted{}
	err := proto.Unmarshal(event, &taskCompleted)
	if err != nil {
		return fmt.Errorf("unmarshal message error: %w", err)
	}

	switch taskCompleted.Payload.(type) {
	case *task_completed.TaskCompleted_V1:
		v1 := taskCompleted.GetV1()
		return c.completeTaskCommand.Handle(ctx, complete_task_command.CompleteTaskData{
			PublicID:    uuid.FromStringOrNil(v1.GetPublicId()),
			CompletedAt: time.Now(), // TODO: use v1.GetCompletedAt() or another event field
		})
	}

	return nil
}
