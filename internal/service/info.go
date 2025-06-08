package service

import (
	"avito_tech/internal/storage"
	"context"
	"fmt"
)

type InfoResponse struct {
	Balance uint64              `json:"balance"`
	Items   []string            `json:"items"`
	SendTo  []storage.CoinsInfo `json:"send_to"`
	GetFrom []storage.CoinsInfo `json:"get_from"`
}
type Info struct {
	repository *storage.UsersDB
}

const (
	SendCoins = "send"
	GetCoins  = "get"
)

func NewInfo(repository *storage.UsersDB) *Info {
	return &Info{repository: repository}
}

func (i *Info) GetUserInfo(ctx context.Context, userId string) (*InfoResponse, error) {
	balance, err := i.getBalance(ctx, userId)
	if err != nil {
		return nil, err
	}

	items, err := i.getItems(ctx, userId)
	if err != nil {
		return nil, err
	}

	sendTo, err := i.getSendTo(ctx, userId, SendCoins)
	if err != nil {
		return nil, err
	}

	getFrom, err := i.getGetFrom(ctx, userId, GetCoins)
	if err != nil {
		return nil, err
	}

	return &InfoResponse{
		Balance: uint64(balance),
		Items:   items,
		SendTo:  sendTo,
		GetFrom: getFrom,
	}, nil

}

func (i *Info) getItems(ctx context.Context, userID string) ([]string, error) {
	items, err := i.repository.GetUserItems(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get items: %w", err)
	}

	return items, nil
}

func (i *Info) getBalance(ctx context.Context, userID string) (int, error) {
	balance, err := i.repository.GetBalance(ctx, userID)
	if err != nil {
		return 0, fmt.Errorf("get balance error: %v", err)
	}

	return balance, nil
}

func (i *Info) getSendTo(ctx context.Context, userID string, send string) ([]storage.CoinsInfo, error) {
	sendTo, err := i.repository.GetCoinsHistory(ctx, userID, send)
	if err != nil {
		return nil, fmt.Errorf("get send to error: %v", err)
	}

	return sendTo, nil
}

func (i *Info) getGetFrom(ctx context.Context, userID string, get string) ([]storage.CoinsInfo, error) {
	getFrom, err := i.repository.GetCoinsHistory(ctx, userID, get)
	if err != nil {
		return nil, fmt.Errorf("get get from error: %v", err)
	}

	return getFrom, nil
}
