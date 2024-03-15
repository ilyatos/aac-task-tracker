package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"
	uuid "github.com/satori/go.uuid"
)

type TaskRepository struct {
	db *sqlx.DB
}

type Task struct {
	PublicID    uuid.UUID  `db:"public_id"`
	Description string     `db:"description"`
	CreditPrice *uint64    `db:"credit_price"`
	DebitPrice  *uint64    `db:"debit_price"`
	CompletedAt *time.Time `db:"completed_at"`
}

func NewTaskRepository(db *sqlx.DB) *TaskRepository {
	return &TaskRepository{db: db}
}

func (r *TaskRepository) Create(ctx context.Context, task Task) error {
	_, err := r.db.NamedExecContext(
		ctx,
		`INSERT INTO task (public_id, description) VALUES (:public_id, :description)`,
		task,
	)
	return err
}

func (r *TaskRepository) GetByPublicID(ctx context.Context, publicID uuid.UUID) (Task, error) {
	var task Task
	err := r.db.GetContext(ctx, &task, `SELECT public_id, description, credit_price, debit_price, completed_at FROM task WHERE public_id = $1`, publicID)
	return task, err
}

func (r *TaskRepository) Update(ctx context.Context, task Task) error {
	res, err := r.db.NamedExecContext(
		ctx,
		`UPDATE task SET description = :description, credit_price = :credit_price, debit_price = :debit_price, completed_at = :completed_at WHERE public_id = :public_id`,
		task,
	)
	if err != nil {
		return err
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return fmt.Errorf("task not found")
	}

	return nil
}
