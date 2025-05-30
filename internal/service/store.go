package service

import (
	"avito_tech/internal/storage"
	"context"
	"errors"
	"strings"
)

type BuyResponse struct {
	Status string `json:"status"`
}

const (
	SUCCESS = "SUCCESS"
)

type StoreService struct {
	repository *storage.UsersDB
}

func NewStoreService(repository *storage.UsersDB) *StoreService {
	return &StoreService{repository: repository}
}

func (s *StoreService) BuyItem(ctx context.Context, userId string, itemName string) (*BuyResponse, error) {
	strings.TrimSpace(itemName)

	itemID, err := s.repository.GetItemID(ctx, itemName)
	if err != nil {
		return nil, err
	}

	balance, err := s.repository.GetBalance(ctx, userId)
	if err != nil {
		return nil, err
	}

	price, err := s.repository.GetPrices(ctx, itemID)
	if err != nil {
		return nil, err
	}

	if price > balance {
		return nil, errors.New("not enought coins")
	}

	err = s.repository.BuyItem(ctx, userId, itemID)
	if err != nil {
		return nil, err
	}

	newBalance := balance - price
	err = s.repository.UpdateCoins(ctx, userId, newBalance)
	if err != nil {
		return nil, err
	}

	return &BuyResponse{
		Status: SUCCESS,
	}, nil
}
