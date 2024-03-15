package create_task

import (
	"context"
	"fmt"
	"math/rand"

	uuid "github.com/satori/go.uuid"
	"github.com/segmentio/kafka-go"
	"google.golang.org/protobuf/proto"

	"github.com/ilyatos/aac-task-tracker/billing/internal/broker"
	"github.com/ilyatos/aac-task-tracker/schema_registry/pkg/meta"
	"github.com/ilyatos/aac-task-tracker/schema_registry/pkg/prices/prices_created"

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

	err = producePricesCreatedEvent(ctx, newTask)
	if err != nil {
		return fmt.Errorf("produce prices created event error: %w", err)
	}

	return nil
}

func getCreditPrice() uint64 {
	return uint64(rand.Intn(40-20) + 20)
}

func getDebitPrice() uint64 {
	return uint64(rand.Intn(20-10) + 10)
}

func producePricesCreatedEvent(ctx context.Context, task repository.Task) error {
	transactionAdded, err := proto.Marshal(&prices_created.PricesCreated{
		Header: &meta.Header{Producer: "billing.create_task"},
		Payload: &prices_created.PricesCreated_V1{
			V1: &prices_created.V1{
				TaskPublicId: task.PublicID.String(),
				Credit:       task.CreditPrice,
				Debit:        task.DebitPrice,
			},
		},
	})
	if err != nil {
		return err
	}

	err = broker.Produce(ctx, "prices-stream", kafka.Message{
		Key:   []byte("PricesCreated"),
		Value: transactionAdded,
	})
	if err != nil {
		return err
	}

	return nil
}
