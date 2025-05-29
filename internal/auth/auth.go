package auth

import (
	"avito_tech/internal/storage"
	"context"
	"fmt"
	"github.com/jackc/pgx/v5"
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

func (a *AuthorizationServiceImpl) Login(ctx context.Context, username string, password string) (*AuthResponse, error) {
	userUUID, storeHash, err := a.userRepository.GetUser(ctx, username)
	if err != nil && !errors.Is(err, pgx.ErrNoRows) {
		return nil, fmt.Errorf("failed to get user %w", err)
	}

	if userUUID == "" {
		hashPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
		if err != nil {
			return nil, fmt.Errorf("failed to hash password: %w", err)
		}

		userUUID, err = a.userRepository.CreateUser(ctx, username, string(hashPassword))
		if err != nil {
			return nil, fmt.Errorf("failed to create user: %w", err)
		}

	} else {
		if bcrypt.CompareHashAndPassword([]byte(storeHash), []byte(password)) != nil {
			return nil, errors.New("invalid password")
		}
	}

	accessToken, err := a.jwtProvider.GenerateAccessToken(userUUID)
	if err != nil {
		return nil, fmt.Errorf("failed to generate access token: %w", err)
	}

	response := &AuthResponse{
		Token: accessToken,
	}

	return response, nil
}
