package repository

import (
	"context"
	"time"

	"github.com/jmoiron/sqlx"
	uuid "github.com/satori/go.uuid"

	time_helper "github.com/ilyatos/aac-task-tracker/billing/internal/app/helper/time"
)

type BillingCycleRepository struct {
	db *sqlx.DB
}

func NewBillingCycleRepository(db *sqlx.DB) *BillingCycleRepository {
	return &BillingCycleRepository{db: db}
}

type LastOpenBillingCycle struct {
	ID      uint64    `db:"id"`
	UserID  uuid.UUID `db:"user_public_id"`
	StartAt time.Time `db:"start_at"`
	EndAt   time.Time `db:"end_at"`
}

func (r *BillingCycleRepository) Create(ctx context.Context, userPublicID uuid.UUID, time time.Time) error {
	_, err := r.db.ExecContext(
		ctx,
		`
			INSERT INTO billing_cycle (user_public_id, start_at, end_at) 
			VALUES ($1, $2, $3) 
			ON CONFLICT (user_public_id, start_at, end_at) DO NOTHING
		`,
		userPublicID,
		time_helper.StartOfDay(time),
		time_helper.EndOfDay(time),
	)
	return err
}

func (r *BillingCycleRepository) CreateTx(ctx context.Context, userPublicID uuid.UUID, startAt, endAt time.Time, tx *sqlx.Tx) error {
	_, err := tx.ExecContext(
		ctx,
		`
			INSERT INTO billing_cycle (user_public_id, start_at, end_at) 
			VALUES ($1, $2, $3) 
			ON CONFLICT (user_public_id, start_at, end_at) DO NOTHING
		`,
		userPublicID,
		startAt,
		endAt,
	)
	return err
}

func (r *BillingCycleRepository) CloseTx(ctx context.Context, id uint64, tx *sqlx.Tx) error {
	_, err := tx.ExecContext(
		ctx,
		`
			UPDATE billing_cycle
			SET is_closed = true
			WHERE id = $1
		`,
		id,
	)
	return err
}

func (r *BillingCycleRepository) GetLastOpenTx(ctx context.Context, userPublicID uuid.UUID, tx *sqlx.Tx) (LastOpenBillingCycle, error) {
	billingCycle := LastOpenBillingCycle{}

	err := tx.GetContext(
		ctx,
		&billingCycle,
		`
			SELECT id, user_public_id, start_at, end_at
			FROM billing_cycle
			WHERE user_public_id = $1 AND is_closed = false
			ORDER BY start_at DESC
			LIMIT 1
			FOR UPDATE
		`,
		userPublicID,
	)
	if err != nil {
		return billingCycle, err
	}

	return billingCycle, nil
}
