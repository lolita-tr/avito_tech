package auth

import (
	"avito_tech/internal/storage"
	"context"
	"github.com/pkg/errors"
	"golang.org/x/crypto/bcrypt"
)

type AuthorizationServiceImpl struct {
	userRepository *storage.UsersDB
	jwtProvider    *JwtProvider
}

func NewAuthorizationService(userRepository *storage.UsersDB, provider *JwtProvider) *AuthorizationServiceImpl {
	return &AuthorizationServiceImpl{
		userRepository: userRepository,
		jwtProvider:    provider,
	}
}

func (a *AuthorizationServiceImpl) Login(ctx context.Context, username string, password string) (string, error) {
	userUUID, storeHash, err := a.userRepository.GetUser(ctx, username)
	if err != nil {
		return "", err
	}

	if userUUID == "" {
		hashPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

		err = a.userRepository.CreateUser(ctx, username, string(hashPassword))
		if err != nil {
			return "", err
		}
	}

	if bcrypt.CompareHashAndPassword([]byte(storeHash), []byte(password)) != nil {
		return "", errors.New("invalid password")
	}

	accessToken, err := a.jwtProvider.GenerateAccessToken(userUUID)
	if err != nil {
		return "", err
	}

	return accessToken, nil
}
