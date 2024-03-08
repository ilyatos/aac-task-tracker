package singup

import (
	"context"
	"encoding/json"
	"fmt"

	uuid "github.com/satori/go.uuid"
	"github.com/segmentio/kafka-go"

	"github.com/ilyatos/aac-task-tracker/auth/internal/app/helper/password"
	"github.com/ilyatos/aac-task-tracker/auth/internal/app/model/user"
	"github.com/ilyatos/aac-task-tracker/auth/internal/app/repository"
	"github.com/ilyatos/aac-task-tracker/auth/internal/broker"
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

type userCreatedEvent struct {
	PublicID uuid.UUID `json:"public_id"`
	Name     string    `json:"name"`
	Email    string    `json:"email"`
	Role     user.Role `json:"role"`
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
	err = produceUserCreatedEvent(ctx, userCreatedEvent{
		PublicID: newUser.PublicID,
		Name:     newUser.Name,
		Email:    newUser.Email,
		Role:     newUser.Role,
	})
	if err != nil {
		return uuid.UUID{}, fmt.Errorf("produce user created event error: %w", err)
	}

	return newUser.PublicID, nil
}

func produceUserCreatedEvent(ctx context.Context, userCreatedEvent userCreatedEvent) error {
	userCreated, err := json.Marshal(userCreatedEvent)
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
