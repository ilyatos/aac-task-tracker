package change_role

import (
	"context"
	"encoding/json"
	"fmt"

	uuid "github.com/satori/go.uuid"
	"github.com/segmentio/kafka-go"

	user_model "github.com/ilyatos/aac-task-tracker/auth/internal/app/model/user"
	"github.com/ilyatos/aac-task-tracker/auth/internal/app/repository"
	"github.com/ilyatos/aac-task-tracker/auth/internal/broker"
)

type Command struct {
	userRepository *repository.UserRepository
}

func New(userRepository *repository.UserRepository) *Command {
	return &Command{userRepository: userRepository}
}

type userRoleChangedEvent struct {
	PublicID uuid.UUID       `json:"public_id"`
	Role     user_model.Role `json:"role"`
}

type userUpdatedEvent struct {
	PublicID uuid.UUID       `json:"public_id"`
	Name     string          `json:"name"`
	Email    string          `json:"email"`
	Role     user_model.Role `json:"role"`
}

func (cr *Command) Handle(ctx context.Context, publicID uuid.UUID, role user_model.Role) error {
	user, err := cr.userRepository.FindByPublicID(ctx, publicID)
	if err != nil {
		return fmt.Errorf("update role error: %w", err)
	}

	user.Role = role

	err = cr.userRepository.Update(ctx, *user)
	if err != nil {
		return fmt.Errorf("update role error: %w", err)
	}

	// TODO: outbox pattern
	err = produceUserUpdatedEvent(ctx, userUpdatedEvent{
		PublicID: publicID,
		Name:     user.Name,
		Email:    user.Email,
		Role:     user.Role,
	})
	if err != nil {
		return fmt.Errorf("produce user role changed event error: %w", err)
	}

	// TODO: outbox pattern
	err = produceUserRoleChangedEvent(ctx, userRoleChangedEvent{PublicID: publicID, Role: role})
	if err != nil {
		return fmt.Errorf("produce user role changed event error: %w", err)
	}

	return nil
}

func produceUserUpdatedEvent(ctx context.Context, userUpdatedEvent userUpdatedEvent) error {
	userUpdated, err := json.Marshal(userUpdatedEvent)
	if err != nil {
		return err
	}

	err = broker.Produce(ctx, "users-stream", kafka.Message{
		Key:   []byte("UserUpdated"),
		Value: userUpdated,
	})
	if err != nil {
		return err
	}

	return nil
}

func produceUserRoleChangedEvent(ctx context.Context, userRoleChangedEvent userRoleChangedEvent) error {
	userChangedRole, err := json.Marshal(userRoleChangedEvent)
	if err != nil {
		return err
	}

	err = broker.Produce(ctx, "users", kafka.Message{
		Key:   []byte("UserRoleChanged"),
		Value: userChangedRole,
	})
	if err != nil {
		return err
	}

	return nil
}
