package repository

import (
	"context"
	"fmt"

	"github.com/jmoiron/sqlx"
	uuid "github.com/satori/go.uuid"

	user_model "github.com/ilyatos/aac-task-tracker/billing/internal/app/model/user"
)

type UserRepository struct {
	db *sqlx.DB
}

type User struct {
	PublicID uuid.UUID       `db:"public_id"`
	Name     string          `db:"name"`
	Email    string          `db:"email"`
	Role     user_model.Role `db:"role"`
	Balance  int64           `db:"balance"`
}

func NewUserRepository(db *sqlx.DB) *UserRepository {
	return &UserRepository{db: db}

}

func (r *UserRepository) CreateTx(ctx context.Context, user User, tx *sqlx.Tx) error {
	_, err := tx.NamedExecContext(
		ctx,
		`INSERT INTO "user" (public_id, name, email, role, balance) VALUES (:public_id, :name, :email, :role, :balance)`,
		user,
	)

	return err
}

func (r *UserRepository) UpdateTx(ctx context.Context, user User, tx *sqlx.Tx) error {
	res, err := tx.NamedExecContext(
		ctx,
		`UPDATE "user" SET name = :name, email = :email, role = :role, balance = :balance WHERE public_id = :public_id`,
		user,
	)
	if err != nil {
		return err
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return fmt.Errorf("user not found")
	}

	return nil
}

func (r *UserRepository) GetForUpdateTx(ctx context.Context, publicID uuid.UUID, tx *sqlx.Tx) (User, error) {
	user := User{}
	err := tx.GetContext(ctx, &user, `SELECT public_id, name, email, role, balance FROM "user" WHERE public_id = $1 FOR UPDATE`, publicID)
	if err != nil {
		return User{}, err
	}

	return user, nil
}
