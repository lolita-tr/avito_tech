package storage

import (
	"context"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/pkg/errors"
)

type UsersDB struct {
	db *pgxpool.Pool
}

func NewUsersDB(db *pgxpool.Pool) *UsersDB {
	return &UsersDB{db: db}
}

const (
	createUserQuery = `
	INSERT INTO users (login, password) VALUES ($1, $2)`
)

func (ud *UsersDB) CreateUser(ctx context.Context, login string, password string) error {
	_, err := ud.db.Exec(ctx, createUserQuery, login, password)
	if err != nil {
		return errors.Wrap(err, "failed to create user")
	}

	return nil
}
