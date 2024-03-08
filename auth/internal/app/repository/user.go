package repository

import (
	"context"
	"fmt"

	"github.com/jmoiron/sqlx"
	uuid "github.com/satori/go.uuid"

	user_model "github.com/ilyatos/aac-task-tracker/auth/internal/app/model/user"
)

type UserRepository struct {
	db *sqlx.DB
}

type User struct {
	PublicID uuid.UUID       `db:"public_id"`
	Name     string          `db:"name"`
	Email    string          `db:"email"`
	Password string          `db:"password"`
	Role     user_model.Role `db:"role"`
}

func NewUserRepository(db *sqlx.DB) *UserRepository {
	return &UserRepository{db: db}

}

func (r *UserRepository) Create(ctx context.Context, user User) error {
	_, err := r.db.NamedExecContext(
		ctx,
		`INSERT INTO "user" (public_id, name, email, password, role) VALUES (:public_id, :name, :email, :password, :role)`,
		user,
	)
	return err
}

func (r *UserRepository) Update(ctx context.Context, user User) error {
	res, err := r.db.NamedExecContext(
		ctx,
		`UPDATE "user" SET name=:name, email=:email, password=:password, role=:role WHERE public_id=:public_id`,
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

	return err
}

func (r *UserRepository) FindByPublicID(ctx context.Context, publicID uuid.UUID) (*User, error) {
	user := &User{}
	err := r.db.GetContext(
		ctx,
		user,
		`SELECT public_id, name, email, password, role FROM "user" WHERE public_id=$1`,
		publicID,
	)
	return user, err
}

func (r *UserRepository) FindByEmail(ctx context.Context, email string) (*User, error) {
	user := &User{}
	err := r.db.GetContext(
		ctx,
		user,
		`SELECT public_id, name, email, password, role FROM "user" WHERE email=$1`,
		email,
	)
	return user, err
}
