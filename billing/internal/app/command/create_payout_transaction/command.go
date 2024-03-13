package create_payout_transaction

import (
	"context"
	"fmt"
	"log"

	"github.com/jmoiron/sqlx"
	uuid "github.com/satori/go.uuid"

	"github.com/ilyatos/aac-task-tracker/billing/internal/app/command/create_transaction"
	time_helper "github.com/ilyatos/aac-task-tracker/billing/internal/app/helper/time"
	transaction_model "github.com/ilyatos/aac-task-tracker/billing/internal/app/model/transaction"
	"github.com/ilyatos/aac-task-tracker/billing/internal/app/repository"
)

type Command struct {
	db                       *sqlx.DB
	billingCycleRepository   *repository.BillingCycleRepository
	transactionRepository    *repository.TransactionRepository
	createTransactionCommand *create_transaction.Command
}

func New(
	db *sqlx.DB,
	billingCycleRepository *repository.BillingCycleRepository,
	transactionRepository *repository.TransactionRepository,
	createTransactionCommand *create_transaction.Command,
) *Command {
	return &Command{
		db:                       db,
		billingCycleRepository:   billingCycleRepository,
		transactionRepository:    transactionRepository,
		createTransactionCommand: createTransactionCommand,
	}
}

func (c *Command) Handle(ctx context.Context, userPublicID uuid.UUID) error {
	log.Println("create payout transaction")

	tx := c.db.MustBeginTx(ctx, nil)

	lastOpenBillingCycle, err := c.billingCycleRepository.GetLastOpenTx(ctx, userPublicID, tx)
	if err != nil {
		_ = tx.Rollback()
		return fmt.Errorf(" get last open billing cycle error: %w", err)
	}

	amount, err := c.transactionRepository.GetBalanceByBillingCycleIDTx(ctx, lastOpenBillingCycle.ID, tx)
	if err != nil {
		_ = tx.Rollback()
		return fmt.Errorf(" get balance error: %w", err)
	}
	if amount >= 0 {
		err = c.createTransactionCommand.Handle(
			ctx,
			create_transaction.Transaction{
				UserPublicID: userPublicID,
				Type:         transaction_model.Payout,
				Credit:       0,
				Debit:        amount,
			},
			tx,
		)
		if err != nil {
			_ = tx.Rollback()
			return fmt.Errorf("create transaction error: %w", err)
		}
	}

	err = c.billingCycleRepository.CloseTx(ctx, lastOpenBillingCycle.ID, tx)
	if err != nil {
		_ = tx.Rollback()
		return fmt.Errorf("close billing cycle error: %w", err)
	}

	err = c.billingCycleRepository.CreateTx(ctx, userPublicID, time_helper.NextDay(lastOpenBillingCycle.StartAt), time_helper.NextDay(lastOpenBillingCycle.EndAt), tx)
	if err != nil {
		_ = tx.Rollback()
		return fmt.Errorf("create billing cycle error: %w", err)
	}

	_ = tx.Commit()

	return nil
}
