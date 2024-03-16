package transactions_lifecycle

import (
	"context"
	"fmt"
	"log"
	"os"
	"strings"

	uuid "github.com/satori/go.uuid"
	"github.com/segmentio/kafka-go"
	"google.golang.org/protobuf/proto"

	create_transaction_command "github.com/ilyatos/aac-task-tracker/analytics/internal/app/command/create_transaction"
	transaction_model "github.com/ilyatos/aac-task-tracker/analytics/internal/app/model/transaction"
	transaction_event "github.com/ilyatos/aac-task-tracker/schema_registry/pkg/transaction"
	"github.com/ilyatos/aac-task-tracker/schema_registry/pkg/transaction/transaction_added"
)

type Consumer struct {
	reader                   *kafka.Reader
	createTransactionCommand *create_transaction_command.Command
}

func New(
	createTransactionCommand *create_transaction_command.Command,
) *Consumer {
	r := kafka.NewReader(kafka.ReaderConfig{
		Brokers:  []string{os.Getenv("KAFKA_URL")},
		GroupID:  "analytics-group",
		Topic:    "transactions",
		MaxBytes: 10e6, // 10MB
	})

	return &Consumer{
		reader:                   r,
		createTransactionCommand: createTransactionCommand,
	}
}

func (c *Consumer) Consume() {
	log.Println("start consuming transactions")

	ctx := context.Background()
	for {
		m, err := c.reader.FetchMessage(ctx)
		if err != nil {
			log.Println("failed to commit messages:", err)
			continue
		}

		log.Printf("message at topic/partition/offset %v/%v/%v: %s = %s\n", m.Topic, m.Partition, m.Offset, string(m.Key), string(m.Value))

		switch string(m.Key) {
		case "TransactionAdded":
			err = c.handleTransactionAdded(ctx, m.Value)
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

func (c *Consumer) handleTransactionAdded(ctx context.Context, event []byte) error {
	log.Println("handle transaction added")

	transactionAdded := transaction_added.TransactionAdded{}
	err := proto.Unmarshal(event, &transactionAdded)
	if err != nil {
		return fmt.Errorf("unmarshal message error: %w", err)
	}

	switch transactionAdded.Payload.(type) {
	case *transaction_added.TransactionAdded_V1:
		v1 := transactionAdded.GetV1()
		return c.createTransactionCommand.Handle(ctx, create_transaction_command.CreateTransactionData{
			PublicID:     uuid.FromStringOrNil(v1.GetPublicId()),
			UserPublicID: uuid.FromStringOrNil(v1.GetUserPublicId()),
			Type:         transaction_model.Type(strings.ToLower(transaction_event.Type_name[int32(v1.GetType())])),
			Credit:       v1.GetCredit(),
			Debit:        v1.GetDebit(),
		})
	}

	return nil
}
