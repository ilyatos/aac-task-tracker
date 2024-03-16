package create_transaction

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/jmoiron/sqlx"
	uuid "github.com/satori/go.uuid"
	"github.com/segmentio/kafka-go"
	"google.golang.org/protobuf/proto"

	transaction_model "github.com/ilyatos/aac-task-tracker/billing/internal/app/model/transaction"
	"github.com/ilyatos/aac-task-tracker/billing/internal/app/repository"
	"github.com/ilyatos/aac-task-tracker/billing/internal/broker"
	"github.com/ilyatos/aac-task-tracker/schema_registry/pkg/meta"
	transaction_event "github.com/ilyatos/aac-task-tracker/schema_registry/pkg/transaction"
	"github.com/ilyatos/aac-task-tracker/schema_registry/pkg/transaction/transaction_added"
)

type Command struct {
	userRepository         *repository.UserRepository
	transactionRepository  *repository.TransactionRepository
	billingCycleRepository *repository.BillingCycleRepository
}

func New(
	userRepository *repository.UserRepository,
	billingCycleRepository *repository.BillingCycleRepository,
	transactionRepository *repository.TransactionRepository,
) *Command {
	return &Command{
		userRepository:         userRepository,
		billingCycleRepository: billingCycleRepository,
		transactionRepository:  transactionRepository,
	}
}

type Transaction struct {
	UserPublicID uuid.UUID
	Type         transaction_model.Type
	Credit       uint64
	Debit        uint64
}

func (c *Command) Handle(ctx context.Context, transaction Transaction, tx *sqlx.Tx) error {
	log.Println("create transaction")

	lastOpenBillingCycle, err := c.billingCycleRepository.GetLastOpenTx(ctx, transaction.UserPublicID, tx)
	if err != nil {
		_ = tx.Rollback()
		return fmt.Errorf("get last open billing cycle error: %w", err)
	}

	user, err := c.userRepository.GetForUpdateTx(ctx, transaction.UserPublicID, tx)
	if err != nil {
		return fmt.Errorf("get last open billing cycle error: %w", err)
	}

	user.Balance += int64(transaction.Credit - transaction.Debit)

	err = c.userRepository.UpdateTx(ctx, user, tx)
	if err != nil {
		return fmt.Errorf("update user error: %w", err)
	}

	log.Println("user balance updated")

	newTransaction := repository.Transaction{
		PublicID:       uuid.NewV4(),
		BillingCycleID: lastOpenBillingCycle.ID,
		UserPublicID:   transaction.UserPublicID,
		Type:           transaction.Type,
		Credit:         transaction.Credit,
		Debit:          transaction.Debit,
	}
	err = c.transactionRepository.CreateTx(ctx, newTransaction, tx)
	if err != nil {
		return fmt.Errorf("create transaction error: %w", err)
	}

	log.Println("transaction created")

	// TODO: outbox pattern
	err = produceTransactionCreatedEvent(ctx, newTransaction)
	if err != nil {
		return fmt.Errorf("produce transaction created event error: %w", err)
	}

	return nil
}

func produceTransactionCreatedEvent(ctx context.Context, transaction repository.Transaction) error {
	transactionAdded, err := proto.Marshal(&transaction_added.TransactionAdded{
		Header: &meta.Header{Producer: "billing.create_transaction"},
		Payload: &transaction_added.TransactionAdded_V1{
			V1: &transaction_added.V1{
				PublicId:     transaction.PublicID.String(),
				UserPublicId: transaction.UserPublicID.String(),
				Type:         transaction_event.Type(transaction_event.Type_value[strings.ToUpper(string(transaction.Type))]),
				Credit:       transaction.Credit,
				Debit:        transaction.Debit,
			},
		},
	})
	if err != nil {
		return err
	}

	err = broker.Produce(ctx, "transactions", kafka.Message{
		Key:   []byte("TransactionAdded"),
		Value: transactionAdded,
	})
	if err != nil {
		return err
	}

	return nil
}
