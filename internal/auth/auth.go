package auth

import (
	"avito_tech/internal/storage"
	"context"
)

type AuthorizationServiceImpl struct {
	userRepository *storage.UsersDB
}

func NewAuthorizationService(userRepository *storage.UsersDB) *AuthorizationServiceImpl {
	return &AuthorizationServiceImpl{userRepository: userRepository}
}

func (a *AuthorizationServiceImpl) Login(ctx context.Context, username string, password string) (string, error) {

}
