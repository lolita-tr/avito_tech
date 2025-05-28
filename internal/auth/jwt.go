package auth

import (
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"github.com/pkg/errors"
	"os"
	"strings"
	"time"
)

type JwtProvider struct {
	jwtKey        []byte
	accessExpity  time.Duration
	signingMethod jwt.SigningMethod
}

func NewJwtProvider() *JwtProvider {
	jwtSecret := os.Getenv("JWT_SECRET")

	return &JwtProvider{
		jwtKey:        []byte(jwtSecret),
		accessExpity:  5 * time.Minute,
		signingMethod: jwt.SigningMethodHS256,
	}
}

func (p *JwtProvider) GenerateAccessToken(userUUID string) (string, error) {
	accessToken := jwt.NewWithClaims(p.signingMethod, jwt.MapClaims{
		"user_id": userUUID,
		"exp":     time.Now().Add(p.accessExpity).Unix(),
	})

	accessTokenString, err := accessToken.SignedString(p.jwtKey)
	if err != nil {
		return "", errors.Wrap(err, "failed to sign access token")
	}

	return accessTokenString, nil
}

func (p *JwtProvider) ValidateAccessToken(accessTokenString string) (string, error) {
	tokenString := strings.TrimSpace(accessTokenString)
	if tokenString == "" {
		return "", errors.New("empty access token")
	}

	token, err := jwt.Parse(accessTokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return p.jwtKey, nil
	})

	if err != nil {
		return "", errors.Wrap(err, "failed to parse access token")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		return "", errors.New("invalid token claims")
	}

	userUUID, ok := claims["user_id"].(string)
	if !ok || userUUID == "" {
		return "", errors.New("user_id not found in token")
	}

	return userUUID, nil
}
