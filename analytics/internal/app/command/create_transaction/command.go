package create_transaction

import (
	"context"
	"fmt"
	"log"

	"github.com/ilyatos/aac-task-tracker/analytics/internal/app/repository"
)

type Command struct {
	transactionRepository *repository.TransactionRepository
}

func New(
	transactionRepository *repository.TransactionRepository,
) *Command {
	return &Command{
		transactionRepository: transactionRepository,
	}
}

type CreateTransactionData repository.Transaction

func (c *Command) Handle(ctx context.Context, transaction CreateTransactionData) error {
	log.Println("create transaction")

	err := c.transactionRepository.Create(ctx, repository.Transaction(transaction))
	if err != nil {
		return fmt.Errorf("create transaction error: %w", err)
	}

	return nil
}
