package create_user

import (
	"context"
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"
	uuid "github.com/satori/go.uuid"

	time_helper "github.com/ilyatos/aac-task-tracker/billing/internal/app/helper/time"
	user_model "github.com/ilyatos/aac-task-tracker/billing/internal/app/model/user"
	"github.com/ilyatos/aac-task-tracker/billing/internal/app/repository"
)

type Command struct {
	db                     *sqlx.DB
	userRepository         *repository.UserRepository
	billingCycleRepository *repository.BillingCycleRepository
}

func New(db *sqlx.DB, userRepository *repository.UserRepository) *Command {
	return &Command{db: db, userRepository: userRepository}
}

type CreateUserData struct {
	PublicID uuid.UUID
	Name     string
	Email    string
	Role     user_model.Role
}

func (c *Command) Handle(ctx context.Context, createUserData CreateUserData) error {
	tx := c.db.MustBeginTx(ctx, nil)

	err := c.userRepository.CreateTx(ctx, repository.User{
		PublicID: createUserData.PublicID,
		Name:     createUserData.Name,
		Email:    createUserData.Email,
		Role:     createUserData.Role,
		Balance:  0,
	}, tx)
	if err != nil {
		_ = tx.Rollback()
		return fmt.Errorf("create user error: %w", err)
	}

	now := time.Now()
	err = c.billingCycleRepository.CreateTx(ctx, createUserData.PublicID, time_helper.StartOfDay(now), time_helper.EndOfDay(now), tx)
	if err != nil {
		_ = tx.Rollback()
		return fmt.Errorf("create billing cycle error: %w", err)
	}

	err = tx.Commit()
	if err != nil {
		return fmt.Errorf("commit transaction error: %w", err)
	}

	return nil
}
