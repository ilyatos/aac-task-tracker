package repository

import (
	"context"

	"github.com/jmoiron/sqlx"
	uuid "github.com/satori/go.uuid"

	transaction_model "github.com/ilyatos/aac-task-tracker/analytics/internal/app/model/transaction"
)

type TransactionRepository struct {
	db *sqlx.DB
}

func NewTransactionRepository(db *sqlx.DB) *TransactionRepository {
	return &TransactionRepository{db: db}
}

type Transaction struct {
	PublicID     uuid.UUID              `db:"public_id"`
	UserPublicID uuid.UUID              `db:"user_public_id"`
	Type         transaction_model.Type `db:"type"`
	Credit       uint64                 `db:"credit"`
	Debit        uint64                 `db:"debit"`
}

func (r *TransactionRepository) Create(ctx context.Context, transaction Transaction) error {
	_, err := r.db.NamedExecContext(
		ctx,
		`
			INSERT INTO transaction (public_id, user_public_id, type, credit, debit)
			VALUES (:public_id, :user_public_id, :type, :credit, :debit)
		`,
		transaction,
	)

	return err
}
