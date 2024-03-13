package update_user

import (
	"context"
	"fmt"

	"github.com/jmoiron/sqlx"
	uuid "github.com/satori/go.uuid"

	user_model "github.com/ilyatos/aac-task-tracker/billing/internal/app/model/user"
	"github.com/ilyatos/aac-task-tracker/billing/internal/app/repository"
)

type Command struct {
	db             *sqlx.DB
	userRepository *repository.UserRepository
}

func New(db *sqlx.DB, userRepository *repository.UserRepository) *Command {
	return &Command{
		db:             db,
		userRepository: userRepository,
	}
}

type UpdateUserData struct {
	PublicID uuid.UUID
	Name     string
	Email    string
	Role     user_model.Role
}

func (c *Command) Handle(ctx context.Context, updateUserData UpdateUserData) error {
	tx := c.db.MustBeginTx(ctx, nil)

	user, err := c.userRepository.GetForUpdateTx(ctx, updateUserData.PublicID, tx)
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("get user for update error: %w", err)
	}

	user.Name = updateUserData.Name
	user.Email = updateUserData.Email
	user.Role = updateUserData.Role

	err = c.userRepository.UpdateTx(ctx, user, tx)
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("create user error: %w", err)
	}

	if err = tx.Commit(); err != nil {
		return fmt.Errorf("commit tx error: %w", err)
	}

	return nil
}
