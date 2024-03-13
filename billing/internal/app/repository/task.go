package repository

import (
	"context"
	"fmt"

	"github.com/jmoiron/sqlx"
	uuid "github.com/satori/go.uuid"
)

type TaskRepository struct {
	db *sqlx.DB
}

type Task struct {
	PublicID     uuid.UUID `db:"public_id"`
	UserPublicID uuid.UUID `db:"user_public_id"`
	Description  string    `db:"description"`
	CreditPrice  uint64    `db:"credit_price"`
	DebitPrice   uint64    `db:"debit_price"`
}

func NewTaskRepository(db *sqlx.DB) *TaskRepository {
	return &TaskRepository{db: db}
}

func (r *TaskRepository) Create(ctx context.Context, task Task) error {
	_, err := r.db.NamedExecContext(
		ctx,
		`INSERT INTO task (public_id, user_public_id, description, credit_price, debit_price) VALUES (:public_id, :user_public_id, :description, :credit_price, :debit_price)`,
		task,
	)
	return err
}

func (r *TaskRepository) Update(ctx context.Context, task Task) error {
	res, err := r.db.NamedExecContext(
		ctx,
		`UPDATE task SET user_public_id = :user_public_id, description = :description WHERE public_id = :public_id`,
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

func (r *TaskRepository) UpdateTx(ctx context.Context, task Task, tx *sqlx.Tx) error {
	res, err := tx.NamedExecContext(
		ctx,
		`UPDATE task SET user_public_id = :user_public_id, description = :description WHERE public_id = :public_id`,
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

func (r *TaskRepository) GetForUpdateTx(ctx context.Context, publicID uuid.UUID, tx *sqlx.Tx) (Task, error) {
	task := Task{}
	err := tx.GetContext(
		ctx,
		&task,
		`SELECT public_id, user_public_id, description, credit_price, debit_price FROM task WHERE public_id = $1 FOR UPDATE`,
		publicID,
	)
	return task, err
}
