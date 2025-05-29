package storage

import (
	"context"
	"errors"
	"fmt"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
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
	SELECT id, password FROM users WHERE login = $1`
)

func (ud *UsersDB) CreateUser(ctx context.Context, login, password string) (string, error) {
	var userID string
	err := ud.db.QueryRow(ctx,
		"INSERT INTO users (login, password) VALUES ($1, $2) RETURNING id",
		login, password,
	).Scan(&userID)

	if err != nil {
		return "", fmt.Errorf("ошибка создания пользователя: %w", err)
	}
	return userID, nil
}

func (ud *UsersDB) GetUser(ctx context.Context, login string) (string, string, error) {
	var userID, storeHash string
	err := ud.db.QueryRow(ctx,
		"SELECT id, password FROM users WHERE login = $1",
		login,
	).Scan(&userID, &storeHash)

	if errors.Is(err, pgx.ErrNoRows) {
		return "", "", nil
	}
	if err != nil {
		return "", "", fmt.Errorf("ошибка при запросе пользователя: %w", err)
	}
	return userID, storeHash, nil
}
