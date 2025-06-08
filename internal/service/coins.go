package service

import (
	"avito_tech/internal/storage"
	"context"
	"fmt"
)

type SendCoinsRequest struct {
	ToUser string `json:"to_user"`
	Amount int    `json:"amount"`
}

type SendCoinsResponse struct {
	Status string `json:"status"`
}

type CoinsService struct {
	repository *storage.UsersDB
}

func NewCoinsService(db *storage.UsersDB) *CoinsService {
	return &CoinsService{repository: db}
}

func (s *CoinsService) SendCoins(ctx context.Context, userID, to_user string, coins_amount int) (*SendCoinsResponse, error) {
	balance, err := s.repository.GetBalance(ctx, userID)
	if err != nil {
		return nil, err
	}

	if balance < coins_amount {
		return &SendCoinsResponse{
			Status: fmt.Sprintf("failed, you have: %d coins", balance),
		}, nil
	}

	to_userID, _, err := s.repository.GetUser(ctx, to_user)
	if err != nil {
		return nil, err
	}

	balance_2, err := s.repository.GetBalance(ctx, to_userID)
	if err != nil {
		return nil, err
	}

	newBalance_user := balance - coins_amount
	newBalance_2 := balance_2 + coins_amount

	err = s.repository.UpdateCoins(ctx, userID, newBalance_user)
	if err != nil {
		return nil, err
	}

	err = s.repository.UpdateCoins(ctx, to_userID, newBalance_2)
	if err != nil {
		return nil, err
	}

	err = s.repository.SaveCoinsHistory(ctx, userID, to_userID, coins_amount)
	if err != nil {
		return nil, err
	}

	return &SendCoinsResponse{Status: SUCCESS}, nil
}
