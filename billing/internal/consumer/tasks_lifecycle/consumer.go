package tasks_lifecycle

import (
	"context"
	"fmt"
	"log"
	"os"

	uuid "github.com/satori/go.uuid"
	"github.com/segmentio/kafka-go"
	"google.golang.org/protobuf/proto"

	"github.com/ilyatos/aac-task-tracker/billing/internal/app/command/create_credit_transaction"
	"github.com/ilyatos/aac-task-tracker/billing/internal/app/command/create_debit_transaction"
	"github.com/ilyatos/aac-task-tracker/schema_registry/pkg/task/task_added"
	"github.com/ilyatos/aac-task-tracker/schema_registry/pkg/task/task_assigned"
	"github.com/ilyatos/aac-task-tracker/schema_registry/pkg/task/task_completed"
)

type Consumer struct {
	reader                  *kafka.Reader
	createDebitTransaction  *create_debit_transaction.Command
	createCreditTransaction *create_credit_transaction.Command
}

func New(
	createDebitTransaction *create_debit_transaction.Command,
	createCreditTransaction *create_credit_transaction.Command,
) *Consumer {
	r := kafka.NewReader(kafka.ReaderConfig{
		Brokers:  []string{os.Getenv("KAFKA_URL")},
		GroupID:  "billing-group",
		Topic:    "tasks",
		MaxBytes: 10e6, // 10MB
	})

	return &Consumer{
		reader:                  r,
		createDebitTransaction:  createDebitTransaction,
		createCreditTransaction: createCreditTransaction,
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
		case "TaskAssigned":
			err = c.handleTaskAssigned(ctx, m.Value)
		case "TaskAdded":
			err = c.handleTaskAdded(ctx, m.Value)
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

func (c *Consumer) handleTaskAssigned(ctx context.Context, event []byte) error {
	log.Println("handle task assigned")

	taskAssigned := task_assigned.TaskAssigned{}
	err := proto.Unmarshal(event, &taskAssigned)
	if err != nil {
		return fmt.Errorf("unmarshal message error: %w", err)
	}

	switch taskAssigned.Payload.(type) {
	case *task_assigned.TaskAssigned_V1:
		v1 := taskAssigned.GetV1()
		return c.createDebitTransaction.Handle(ctx, uuid.FromStringOrNil(v1.GetPublicId()), uuid.FromStringOrNil(v1.GetUserPublicId()))
	}

	return nil
}

func (c *Consumer) handleTaskAdded(ctx context.Context, event []byte) error {
	log.Println("handle task added")

	taskAdded := task_added.TaskAdded{}
	err := proto.Unmarshal(event, &taskAdded)
	if err != nil {
		return fmt.Errorf("unmarshal message error: %w", err)
	}

	switch taskAdded.Payload.(type) {
	case *task_added.TaskAdded_V1:
		v1 := taskAdded.GetV1()
		return c.createDebitTransaction.Handle(ctx, uuid.FromStringOrNil(v1.GetPublicId()), uuid.FromStringOrNil(v1.GetUserPublicId()))
	}

	return nil
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
		return c.createCreditTransaction.Handle(ctx, uuid.FromStringOrNil(v1.GetPublicId()), uuid.FromStringOrNil(v1.GetUserPublicId()))
	}

	return nil
}
