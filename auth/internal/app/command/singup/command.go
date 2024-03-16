package singup

import (
	"context"
	"fmt"
	"strings"

	uuid "github.com/satori/go.uuid"
	"github.com/segmentio/kafka-go"
	"google.golang.org/protobuf/proto"

	"github.com/ilyatos/aac-task-tracker/auth/internal/app/helper/password"
	"github.com/ilyatos/aac-task-tracker/auth/internal/app/model/user"
	"github.com/ilyatos/aac-task-tracker/auth/internal/app/repository"
	"github.com/ilyatos/aac-task-tracker/auth/internal/broker"
	"github.com/ilyatos/aac-task-tracker/schema_registry/pkg/meta"
	user_event "github.com/ilyatos/aac-task-tracker/schema_registry/pkg/user"
	"github.com/ilyatos/aac-task-tracker/schema_registry/pkg/user/user_created"
)

type Command struct {
	userRepository *repository.UserRepository
	hasher         *password.Hasher
}

func New(userRepository *repository.UserRepository, hasher *password.Hasher) *Command {
	return &Command{userRepository: userRepository, hasher: hasher}
}

type UserSignUpData struct {
	Name     string
	Email    string
	Password string
	Role     user.Role
}

func (s *Command) Handle(ctx context.Context, userSignUpData UserSignUpData) (uuid.UUID, error) {
	passwordHash, err := s.hasher.Hash(userSignUpData.Password)
	if err != nil {
		return uuid.UUID{}, fmt.Errorf("hash password error: %w", err)
	}

	newUser := repository.User{
		PublicID: uuid.NewV4(),
		Name:     userSignUpData.Name,
		Email:    userSignUpData.Email,
		Password: passwordHash,
		Role:     userSignUpData.Role,
	}
	err = s.userRepository.Create(ctx, newUser)
	if err != nil {
		return uuid.UUID{}, fmt.Errorf("create user error: %w", err)
	}

	// TODO: outbox pattern
	err = produceUserCreatedEvent(ctx, newUser)
	if err != nil {
		return uuid.UUID{}, fmt.Errorf("produce user created event error: %w", err)
	}

	return newUser.PublicID, nil
}

func produceUserCreatedEvent(ctx context.Context, user repository.User) error {
	userCreated, err := proto.Marshal(&user_created.UserCreated{
		Header: &meta.Header{
			Producer: "auth.signup",
		},
		Payload: &user_created.UserCreated_V1{
			V1: &user_created.V1{
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
		Key:   []byte("UserCreated"),
		Value: userCreated,
	})
	if err != nil {
		return err
	}

	return nil
}
