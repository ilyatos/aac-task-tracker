package repository

import (
	"context"
	"fmt"

	"github.com/jmoiron/sqlx"
	uuid "github.com/satori/go.uuid"

	task_model "github.com/ilyatos/aac-task-tracker/tracker/internal/app/model/task"
)

type TaskRepository struct {
	db *sqlx.DB
}

type NewTask struct {
	PublicID     uuid.UUID         `db:"public_id"`
	UserPublicID uuid.UUID         `db:"user_public_id"`
	Description  string            `db:"description"`
	Status       task_model.Status `db:"status"`
}

type Task struct {
	PublicID     uuid.UUID         `db:"public_id"`
	UserPublicID uuid.UUID         `db:"user_public_id"`
	Description  string            `db:"description"`
	Status       task_model.Status `db:"status"`
}

func NewTaskRepository(db *sqlx.DB) *TaskRepository {
	return &TaskRepository{db: db}
}

func (r *TaskRepository) Create(ctx context.Context, task NewTask) error {
	_, err := r.db.NamedExecContext(
		ctx,
		`INSERT INTO task (public_id, user_public_id, description, status) VALUES (:public_id, :user_public_id, :description, :status)`,
		task,
	)
	return err
}

func (r *TaskRepository) Update(ctx context.Context, task Task) error {
	res, err := r.db.NamedExecContext(
		ctx,
		`UPDATE task SET user_public_id = :user_public_id, description = :description, status = :status WHERE public_id = :public_id`,
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

func (r *TaskRepository) GetByPublicID(ctx context.Context, publicID uuid.UUID) (*Task, error) {
	task := &Task{}
	err := r.db.GetContext(ctx, task, `SELECT public_id, user_public_id, description, status FROM task WHERE public_id = $1`, publicID)
	return task, err
}

func (r *TaskRepository) GetAllByStatus(ctx context.Context, status task_model.Status) ([]Task, error) {
	var tasks []Task
	err := r.db.SelectContext(ctx, &tasks, `SELECT public_id, user_public_id, description, status FROM task WHERE status = $1`, status)
	return tasks, err
}
