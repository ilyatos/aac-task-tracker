package users_stream

import (
	"context"
	"encoding/json"
	"log"
	"os"

	uuid "github.com/satori/go.uuid"
	"github.com/segmentio/kafka-go"

	create_user_command "github.com/ilyatos/aac-task-tracker/tracker/internal/app/command/create_user"
	update_user_command "github.com/ilyatos/aac-task-tracker/tracker/internal/app/command/update_user"
	"github.com/ilyatos/aac-task-tracker/tracker/internal/app/model/user"
)

type Consumer struct {
	reader            *kafka.Reader
	createUserCommand *create_user_command.Command
	updateUserCommand *update_user_command.Command
}

type userCreatedEvent struct {
	PublicID uuid.UUID `json:"public_id"`
	Name     string    `json:"name"`
	Email    string    `json:"email"`
	Role     user.Role `json:"role"`
}

type userUpdatedEvent struct {
	PublicID uuid.UUID `json:"public_id"`
	Name     string    `json:"name"`
	Email    string    `json:"email"`
	Role     user.Role `json:"role"`
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
			userCreatedEvent := userCreatedEvent{}
			err = json.Unmarshal(m.Value, &userCreatedEvent)
			if err != nil {
				log.Println("failed to unmarshal message:", err)
				continue
			}

			err = c.createUserCommand.Handle(ctx, create_user_command.CreateUserData(userCreatedEvent))
		case "UserUpdated":
			userUpdatedEvent := userUpdatedEvent{}
			err = json.Unmarshal(m.Value, &userUpdatedEvent)
			if err != nil {
				log.Println("failed to unmarshal message:", err)
				continue
			}

			err = c.updateUserCommand.Handle(ctx, update_user_command.UpdateUserData(userUpdatedEvent))
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
