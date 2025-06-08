package service

import (
	"avito_tech/internal/storage"
	"context"
	"errors"
	"fmt"
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
		return nil, fmt.Errorf("get item id failed: %v", err)
	}

	balance, err := s.repository.GetBalance(ctx, userId)
	if err != nil {
		return nil, fmt.Errorf("get balance failed: %v", err)
	}

	price, err := s.repository.GetPrices(ctx, itemID)
	if err != nil {
		return nil, fmt.Errorf("get price failed: %v", err)
	}

	if price > balance {
		return nil, errors.New("not enought coins")
	}

	err = s.repository.BuyItem(ctx, userId, itemID)
	if err != nil {
		return nil, fmt.Errorf("buy item failed: %v", err)
	}

	newBalance := balance - price
	err = s.repository.UpdateCoins(ctx, userId, newBalance)
	if err != nil {
		return nil, fmt.Errorf("update coins failed: %v", err)
	}

	return &BuyResponse{
		Status: SUCCESS,
	}, nil
}
