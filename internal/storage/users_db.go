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

type CoinsInfo struct {
	UserId string
	Coins  uint64
}

func NewUsersDB(db *pgxpool.Pool) *UsersDB {
	return &UsersDB{db: db}
}

const (
	createUserQuery = `
	INSERT INTO users (login, password) VALUES ($1, $2) RETURNING id`

	getUserQuery = `
	SELECT id, password FROM users WHERE login = $1`

	getUserItemsQuery = `
	SELECT item_id FROM purch_history WHERE user_id = $1`

	updateBalanceQuery = `
	INSERT INTO balance (user_id, coins_amount) VALUES ($1, $2)`

	getBalancesQuery = `
	SELECT coins_amount FROM balance WHERE user_id = $1
	ORDER BY date DESC LIMIT 1`

	addItemQuery = `
	INSERT INTO purch_history (user_id, item_id) VALUES ($1, $2)`

	getPricesQuery = `
	SELECT price FROM items WHERE id = $1`

	getItemNameQuery = `
	SELECT name FROM items WHERE id = $1`

	getItemIdQuery = `
	SELECT id FROM items WHERE name = $1`

	saveCoinsHistoryQuery = `
	INSERT INTO coins_history (from_user, to_user, coins_amount) VALUES ($1, $2, $3)`

	getSendCoinsHistoryQuery = `
	SELECT to_user, coins_amount FROM coins_history WHERE from_user = $1`

	getGetCoinsHistoryQuery = `
	SELECT from_user, coins_amount FROM coins_history WHERE to_user = $1`
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

	_, err = ud.db.Exec(ctx, updateBalanceQuery, userID, 1000)

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

func (ud *UsersDB) UpdateCoins(ctx context.Context, userID string, coins int) error {

	_, err := ud.db.Exec(ctx, updateBalanceQuery, userID, coins)

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

func (ud *UsersDB) GetPrices(ctx context.Context, itemID string) (int, error) {
	var prices int
	err := ud.db.QueryRow(ctx, getPricesQuery, itemID).Scan(&prices)

	if err != nil {
		return 0, fmt.Errorf("failed to get prices: %w", err)
	}

	return prices, nil
}

func (ud *UsersDB) GetItemName(ctx context.Context, itemID string) (string, error) {
	var itemName string
	err := ud.db.QueryRow(ctx, getItemNameQuery, itemID).Scan(&itemName)

	if err != nil {
		return "", fmt.Errorf("failed to get item name: %w", err)
	}

	return itemName, nil
}

func (ud *UsersDB) GetItemID(ctx context.Context, itemName string) (string, error) {
	var itemID string
	err := ud.db.QueryRow(ctx, getItemIdQuery, itemName).Scan(&itemID)

	if err != nil {
		return "", fmt.Errorf("failed to get item id: %w", err)
	}

	return itemID, nil
}

func (ud *UsersDB) SaveCoinsHistory(ctx context.Context, userID string, to_userID string, coins_amount int) error {
	_, err := ud.db.Exec(ctx, saveCoinsHistoryQuery, userID, to_userID, coins_amount)

	if err != nil {
		return fmt.Errorf("failed to save coins history: %w", err)
	}
	return nil
}

func (ud *UsersDB) GetUserItems(ctx context.Context, userID string) ([]string, error) {
	var items []string
	rows, err := ud.db.Query(ctx, getUserItemsQuery, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to query items: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var itemID string
		if err = rows.Scan(&itemID); err != nil {
			return nil, fmt.Errorf("failed to scan item id: %w", err)
		}

		items = append(items, itemID)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("failed to get items: %w", err)
	}

	return items, nil
}

func (ud *UsersDB) GetCoinsHistory(ctx context.Context, userID string, option string) ([]CoinsInfo, error) {
	var rows pgx.Rows
	var err error

	switch option {
	case "send":
		rows, err = ud.db.Query(ctx, getSendCoinsHistoryQuery, userID)
	case "get":
		rows, err = ud.db.Query(ctx, getGetCoinsHistoryQuery, userID)
	default:
		return nil, fmt.Errorf("invalid option")
	}

	if err != nil {
		return nil, fmt.Errorf("failed to query coins history: %w", err)
	}
	defer rows.Close()

	var coinsHistory []CoinsInfo
	for rows.Next() {
		var coinsInfo CoinsInfo
		if err = rows.Scan(&coinsInfo.UserId, &coinsInfo.Coins); err != nil {
			return nil, fmt.Errorf("failed to scan coins history: %w", err)
		}
		coinsHistory = append(coinsHistory, coinsInfo)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("failed to get coins history: %w", err)
	}

	return coinsHistory, nil
}
