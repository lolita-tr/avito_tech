package service

import (
	"avito_tech/internal/storage"
	"context"
)

type StoreService struct {
	repository storage.UsersDB
}

func NewStoreService(repository storage.UsersDB) *StoreService {
	return &StoreService{repository: repository}
}

func (s *StoreService) BuyItem(ctx context.Context, userId string, itemName string) error {

}
