package repository

import (
	"context"

	"github.com/jmoiron/sqlx"
	uuid "github.com/satori/go.uuid"
)

type UserRepository struct {
	db *sqlx.DB
}

type User struct {
	PublicID uuid.UUID `db:"public_id"`
	Name     string    `db:"name"`
}

func NewUserRepository(db *sqlx.DB) *UserRepository {
	return &UserRepository{db: db}

}

func (r *UserRepository) Create(ctx context.Context, user User) error {
	_, err := r.db.NamedExecContext(
		ctx,
		`INSERT INTO "user" (public_id, name) VALUES (:public_id, :name)`,
		user,
	)

	return err
}
