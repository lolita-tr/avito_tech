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

	getUserQuery = `
	SELECT id, passsword FROM users WHERE login = $1`
)

func (ud *UsersDB) CreateUser(ctx context.Context, login, password string) error {
	_, err := ud.db.Exec(ctx, createUserQuery, login, password)
	if err != nil {
		return errors.Wrap(err, "failed to create user")
	}

	return nil
}

func (ud *UsersDB) GetUser(ctx context.Context, login string) (string, string, error) {
	var userID, storeHash string
	err := ud.db.QueryRow(ctx, getUserQuery, login).Scan(&userID, &storeHash)
	if err != nil {
		return "", "", errors.Wrap(err, "user not found")
	}

	return userID, storeHash, nil
}
