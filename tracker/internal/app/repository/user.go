package repository

import (
	"context"
	"fmt"

	"github.com/jmoiron/sqlx"
	uuid "github.com/satori/go.uuid"

	user_model "github.com/ilyatos/aac-task-tracker/tracker/internal/app/model/user"
)

type UserRepository struct {
	db *sqlx.DB
}

type User struct {
	PublicID uuid.UUID       `db:"public_id"`
	Name     string          `db:"name"`
	Email    string          `db:"email"`
	Role     user_model.Role `db:"role"`
}

type UserPublicID struct {
	ID       string    `db:"id"`
	PublicID uuid.UUID `db:"public_id"`
}

func NewUserRepository(db *sqlx.DB) *UserRepository {
	return &UserRepository{db: db}

}

func (r *UserRepository) Create(ctx context.Context, user User) error {
	_, err := r.db.NamedExecContext(
		ctx,
		`INSERT INTO "user" (public_id, name, email, role) VALUES (:public_id, :name, :email, :role)`,
		user,
	)

	return err
}

func (r *UserRepository) Update(ctx context.Context, user User) error {
	res, err := r.db.NamedExecContext(
		ctx,
		`UPDATE "user" SET name = :name, email = :email, role = :role WHERE public_id = :public_id`,
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

func (r *UserRepository) GetAllByRole(ctx context.Context, role user_model.Role) ([]UserPublicID, error) {
	users := []UserPublicID{}
	err := r.db.SelectContext(
		ctx,
		&users,
		`SELECT id, public_id FROM "user" WHERE role=$1 `,
		role,
	)

	return users, err
}
