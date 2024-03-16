package create_credit_transaction

import (
	"context"
	"fmt"
	"log"

	"github.com/jmoiron/sqlx"
	uuid "github.com/satori/go.uuid"

	"github.com/ilyatos/aac-task-tracker/billing/internal/app/command/create_transaction"
	transaction_model "github.com/ilyatos/aac-task-tracker/billing/internal/app/model/transaction"
	"github.com/ilyatos/aac-task-tracker/billing/internal/app/repository"
)

type Command struct {
	db                       *sqlx.DB
	taskRepository           *repository.TaskRepository
	createTransactionCommand *create_transaction.Command
}

func New(
	db *sqlx.DB,
	taskRepository *repository.TaskRepository,
	createTransactionCommand *create_transaction.Command,
) *Command {
	return &Command{
		db:                       db,
		taskRepository:           taskRepository,
		createTransactionCommand: createTransactionCommand,
	}
}

func (c *Command) Handle(ctx context.Context, taskPublicID, userPublicID uuid.UUID) error {
	log.Println("create credit transaction")

	tx := c.db.MustBeginTx(ctx, nil)

	task, err := c.taskRepository.GetForUpdateTx(ctx, taskPublicID, tx)
	if err != nil {
		_ = tx.Rollback()
		return fmt.Errorf("get task error: %w", err)
	}

	err = c.createTransactionCommand.Handle(
		ctx,
		create_transaction.Transaction{
			UserPublicID: userPublicID,
			Type:         transaction_model.TaskCompleted,
			Credit:       task.CreditPrice,
			Debit:        0,
		},
		tx,
	)
	if err != nil {
		_ = tx.Rollback()
		return fmt.Errorf("create transaction error: %w", err)
	}

	_ = tx.Commit()

	return nil
}
