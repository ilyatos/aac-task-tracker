package update_task

import (
	"context"
	"fmt"

	"github.com/jmoiron/sqlx"
	uuid "github.com/satori/go.uuid"

	"github.com/ilyatos/aac-task-tracker/billing/internal/app/repository"
)

type Command struct {
	db             *sqlx.DB
	taskRepository *repository.TaskRepository
}

func New(db *sqlx.DB, taskRepository *repository.TaskRepository) *Command {
	return &Command{
		db:             db,
		taskRepository: taskRepository,
	}
}

type UpdateTaskData struct {
	PublicID     uuid.UUID
	UserPublicID uuid.UUID
	Description  string
}

func (c *Command) Handle(ctx context.Context, updateTaskData UpdateTaskData) error {
	tx := c.db.MustBeginTx(ctx, nil)

	task, err := c.taskRepository.GetForUpdateTx(ctx, updateTaskData.PublicID, tx)
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("get task error: %w", err)
	}

	task.UserPublicID = updateTaskData.UserPublicID
	task.Description = updateTaskData.Description

	err = c.taskRepository.UpdateTx(ctx, task, tx)
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("update task error: %w", err)
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("commit tx error: %w", err)
	}

	return nil
}
