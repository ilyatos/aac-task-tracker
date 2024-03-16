package users_stream

import (
	"context"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/ilyatos/aac-task-tracker/schema_registry/pkg/user/user_updated"

	user_event "github.com/ilyatos/aac-task-tracker/schema_registry/pkg/user"
	user_model "github.com/ilyatos/aac-task-tracker/tracker/internal/app/model/user"

	uuid "github.com/satori/go.uuid"

	"github.com/ilyatos/aac-task-tracker/schema_registry/pkg/user/user_created"

	"github.com/segmentio/kafka-go"
	"google.golang.org/protobuf/proto"

	create_user_command "github.com/ilyatos/aac-task-tracker/tracker/internal/app/command/create_user"
	update_user_command "github.com/ilyatos/aac-task-tracker/tracker/internal/app/command/update_user"
)

type Consumer struct {
	reader            *kafka.Reader
	createUserCommand *create_user_command.Command
	updateUserCommand *update_user_command.Command
}

func NewConsumer(createUserCommand *create_user_command.Command, updateUserCommand *update_user_command.Command) *Consumer {
	r := kafka.NewReader(kafka.ReaderConfig{
		Brokers:  []string{os.Getenv("KAFKA_URL")},
		GroupID:  "tracker-group",
		Topic:    "users-stream",
		MaxBytes: 10e6, // 10MB
	})

	return &Consumer{reader: r, createUserCommand: createUserCommand, updateUserCommand: updateUserCommand}
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
		case "UserUpdated":
			err = c.handleUserUpdated(ctx, m.Value)
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
			Email:    v1.Email,
			Role:     user_model.Role(strings.ToLower(user_event.Role_name[int32(v1.Role)])),
		})
	}

	return nil
}

func (c *Consumer) handleUserUpdated(ctx context.Context, event []byte) error {
	userUpdated := user_updated.UserUpdated{}
	err := proto.Unmarshal(event, &userUpdated)
	if err != nil {
		return fmt.Errorf("unmarshal message error: %w", err)
	}

	switch userUpdated.Payload.(type) {
	case *user_updated.UserUpdated_V1:
		v1 := userUpdated.GetV1()
		return c.updateUserCommand.Handle(ctx, update_user_command.UpdateUserData{
			PublicID: uuid.FromStringOrNil(v1.PublicId),
			Name:     v1.Name,
			Email:    v1.Email,
			Role:     user_model.Role(strings.ToLower(user_event.Role_name[int32(v1.Role)])),
		})
	}

	return nil
}
