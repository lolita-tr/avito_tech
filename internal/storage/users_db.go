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
	INSERT INTO users (login, password) VALUES ($1, $2) RETURNING id`

	getUserQuery = `
	SELECT id, password FROM users WHERE login = $1`

	addCoinsQuery = `
	INSERT INTO balance (user_id, coins_amount) VALUES ($1, $2)`

	getBalancesQuery = `
	SELECT coins_amount FROM balance WHERE user_id = $1
	ORDER BY date DESC LIMIT 1`

	addItemQuery = `
	INSERT INTO purch_history (user_id, item_id) VALUES ($1, $2)`
)

func (ud *UsersDB) CreateUser(ctx context.Context, login, password string) (string, error) {
	var userID string
	err := ud.db.QueryRow(ctx,
		createUserQuery,
		login, password,
	).Scan(&userID)

	if err != nil {
		return "", fmt.Errorf("failed to create user: %w", err)
	}

	_, err = ud.db.Exec(ctx, addCoinsQuery, userID, 1000)

	return userID, nil
}

func (ud *UsersDB) GetUser(ctx context.Context, login string) (string, string, error) {
	var userID, storeHash string
	err := ud.db.QueryRow(ctx,
		getUserQuery,
		login,
	).Scan(&userID, &storeHash)

	if errors.Is(err, pgx.ErrNoRows) {
		return "", "", nil
	}
	if err != nil {
		return "", "", fmt.Errorf("failed to get user: %w", err)
	}
	return userID, storeHash, nil
}

func (ud *UsersDB) GetBalance(ctx context.Context, userID string) (int, error) {
	var balance int
	err := ud.db.QueryRow(ctx, getBalancesQuery, userID).Scan(&balance)

	if err != nil {
		return 0, fmt.Errorf("failed to get balance: %w", err)
	}

	return balance, nil
}

func (ud *UsersDB) AddCoins(ctx context.Context, userID string, coins int) error {
	currentBalance, err := ud.GetBalance(ctx, userID)

	if err != nil {
		return fmt.Errorf("failed to get current balance: %w", err)
	}

	newBalance := currentBalance + coins
	_, err = ud.db.Exec(ctx, addCoinsQuery, userID, newBalance)

	if err != nil {
		return fmt.Errorf("failed to add coins: %w", err)
	}

	return nil
}

func (ud *UsersDB) BuyItem(ctx context.Context, userID string, itemID string) error {
	_, err := ud.db.Exec(ctx, addItemQuery, userID, itemID)
	if err != nil {
		return fmt.Errorf("failed to buy item: %w", err)
	}

	return nil
}
