package repository

import (
	"context"

	"github.com/jmoiron/sqlx"
	uuid "github.com/satori/go.uuid"

	transaction_model "github.com/ilyatos/aac-task-tracker/billing/internal/app/model/transaction"
)

type TransactionRepository struct {
	db *sqlx.DB
}

func NewTransactionRepository(db *sqlx.DB) *TransactionRepository {
	return &TransactionRepository{db: db}
}

type Transaction struct {
	PublicID       uuid.UUID              `db:"public_id"`
	BillingCycleID uint64                 `db:"billing_cycle_id"`
	UserPublicID   uuid.UUID              `db:"user_public_id"`
	Type           transaction_model.Type `db:"type"`
	Credit         uint64                 `db:"credit"`
	Debit          uint64                 `db:"debit"`
}

func (r *TransactionRepository) CreateTx(ctx context.Context, transaction Transaction, tx *sqlx.Tx) error {
	_, err := tx.NamedExecContext(
		ctx,
		`
			INSERT INTO transaction (public_id, billing_cycle_id, user_public_id, type, credit, debit)
			VALUES (:public_id, :billing_cycle_id, :user_public_id, :type, :credit, :debit)
		`,
		transaction,
	)

	return err
}

func (r *TransactionRepository) GetBalanceByBillingCycleIDTx(ctx context.Context, billingCycleID uint64, tx *sqlx.Tx) (uint64, error) {
	var balance uint64
	err := tx.GetContext(
		ctx,
		&balance,
		`
			SELECT SUM(credit) - SUM(debit) AS balance
			FROM transaction
			WHERE billing_cycle_id = $1
		`,
		billingCycleID,
	)
	return balance, err
}
