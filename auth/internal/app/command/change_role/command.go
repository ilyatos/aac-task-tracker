package change_role

import (
	"context"
	"fmt"
	"strings"

	uuid "github.com/satori/go.uuid"
	"github.com/segmentio/kafka-go"
	"google.golang.org/protobuf/proto"

	user_model "github.com/ilyatos/aac-task-tracker/auth/internal/app/model/user"
	"github.com/ilyatos/aac-task-tracker/auth/internal/app/repository"
	"github.com/ilyatos/aac-task-tracker/auth/internal/broker"
	"github.com/ilyatos/aac-task-tracker/schema_registry/pkg/meta"
	user_event "github.com/ilyatos/aac-task-tracker/schema_registry/pkg/user"
	"github.com/ilyatos/aac-task-tracker/schema_registry/pkg/user/role_changed"
	"github.com/ilyatos/aac-task-tracker/schema_registry/pkg/user/user_updated"
)

type Command struct {
	userRepository *repository.UserRepository
}

func New(userRepository *repository.UserRepository) *Command {
	return &Command{userRepository: userRepository}
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
	err = produceUserUpdatedEvent(ctx, user)
	if err != nil {
		return fmt.Errorf("produce user role changed event error: %w", err)
	}

	// TODO: outbox pattern
	err = produceUserRoleChangedEvent(ctx, user)
	if err != nil {
		return fmt.Errorf("produce user role changed event error: %w", err)
	}

	return nil
}

func produceUserUpdatedEvent(ctx context.Context, user *repository.User) error {
	userUpdated, err := proto.Marshal(&user_updated.UserUpdated{
		Header: &meta.Header{
			Producer: "auth.change_role",
		},
		Payload: &user_updated.UserUpdated_V1{
			V1: &user_updated.V1{
				PublicId: user.PublicID.String(),
				Name:     user.Name,
				Email:    user.Email,
				Role:     user_event.Role(user_event.Role_value[strings.ToUpper(string(user.Role))]),
			},
		},
	})
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

func produceUserRoleChangedEvent(ctx context.Context, user *repository.User) error {
	userRoleChanged, err := proto.Marshal(&role_changed.UserRoleChanged{
		Header: &meta.Header{
			Producer: "auth.change_role",
		},
		Payload: &role_changed.UserRoleChanged_V1{
			V1: &role_changed.V1{
				PublicId: user.PublicID.String(),
				Role:     user_event.Role(user_event.Role_value[strings.ToUpper(string(user.Role))]),
			},
		},
	})
	if err != nil {
		return err
	}

	err = broker.Produce(ctx, "users", kafka.Message{
		Key:   []byte("UserRoleChanged"),
		Value: userRoleChanged,
	})
	if err != nil {
		return err
	}

	return nil
}
