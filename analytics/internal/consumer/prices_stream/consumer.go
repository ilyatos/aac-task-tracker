package prices_stream

import (
	"context"
	"fmt"
	"log"
	"os"

	uuid "github.com/satori/go.uuid"
	"github.com/segmentio/kafka-go"
	"google.golang.org/protobuf/proto"

	add_task_prices_command "github.com/ilyatos/aac-task-tracker/analytics/internal/app/command/add_task_prices"
	"github.com/ilyatos/aac-task-tracker/schema_registry/pkg/prices/prices_created"
)

type Consumer struct {
	reader        *kafka.Reader
	addTaskPrices *add_task_prices_command.Command
}

func New(
	addTaskPrices *add_task_prices_command.Command,
) *Consumer {
	r := kafka.NewReader(kafka.ReaderConfig{
		Brokers:  []string{os.Getenv("KAFKA_URL")},
		GroupID:  "analytics-group",
		Topic:    "prices-stream",
		MaxBytes: 10e6, // 10MB
	})

	return &Consumer{
		reader:        r,
		addTaskPrices: addTaskPrices,
	}
}

func (c *Consumer) Consume() {
	log.Println("start consuming prices stream")

	ctx := context.Background()
	for {
		m, err := c.reader.FetchMessage(ctx)
		if err != nil {
			log.Println("failed to commit messages:", err)
			continue
		}

		log.Printf("message at topic/partition/offset %v/%v/%v: %s = %s\n", m.Topic, m.Partition, m.Offset, string(m.Key), string(m.Value))

		switch string(m.Key) {
		case "PricesCreated":
			err = c.handlePricesCreated(ctx, m.Value)
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

func (c *Consumer) handlePricesCreated(ctx context.Context, event []byte) error {
	pricesCreated := prices_created.PricesCreated{}
	err := proto.Unmarshal(event, &pricesCreated)
	if err != nil {
		return fmt.Errorf("unmarshal message error: %w", err)
	}

	switch pricesCreated.Payload.(type) {
	case *prices_created.PricesCreated_V1:
		v1 := pricesCreated.GetV1()
		return c.addTaskPrices.Handle(ctx, add_task_prices_command.AddTaskPricesData{
			TaskPublicID: uuid.FromStringOrNil(v1.GetTaskPublicId()),
			CreditPrice:  v1.GetCredit(),
			DebitPrice:   v1.GetDebit(),
		})
	}

	return nil
}
