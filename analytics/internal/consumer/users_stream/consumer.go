package users_stream

import (
	"context"
	"fmt"
	"log"
	"os"

	uuid "github.com/satori/go.uuid"
	"github.com/segmentio/kafka-go"
	"google.golang.org/protobuf/proto"

	create_user_command "github.com/ilyatos/aac-task-tracker/analytics/internal/app/command/create_user"
	"github.com/ilyatos/aac-task-tracker/schema_registry/pkg/user/user_created"
)

type Consumer struct {
	reader            *kafka.Reader
	createUserCommand *create_user_command.Command
}

func New(createUserCommand *create_user_command.Command) *Consumer {
	r := kafka.NewReader(kafka.ReaderConfig{
		Brokers:  []string{os.Getenv("KAFKA_URL")},
		GroupID:  "analytics-group",
		Topic:    "users-stream",
		MaxBytes: 10e6, // 10MB
	})

	return &Consumer{reader: r, createUserCommand: createUserCommand}
}

func (c *Consumer) Consume() {
	log.Println("start consuming users stream")

	ctx := context.Background()
	for {
		m, err := c.reader.FetchMessage(ctx)
		if err != nil {
			log.Println("failed to commit messages:", err)
			continue
		}

		log.Printf("message at topic/partition/offset %v/%v/%v: %s = %s\n", m.Topic, m.Partition, m.Offset, string(m.Key), string(m.Value))

		switch string(m.Key) {
		case "UserCreated":
			err = c.handleUserCreated(ctx, m.Value)
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

func (c *Consumer) handleUserCreated(ctx context.Context, event []byte) error {
	userCreated := user_created.UserCreated{}
	err := proto.Unmarshal(event, &userCreated)
	if err != nil {
		return fmt.Errorf("unmarshal message error: %w", err)
	}

	switch userCreated.Payload.(type) {
	case *user_created.UserCreated_V1:
		v1 := userCreated.GetV1()
		return c.createUserCommand.Handle(ctx, create_user_command.CreateUserData{
			PublicID: uuid.FromStringOrNil(v1.PublicId),
			Name:     v1.Name,
		})
	}

	return nil
}
