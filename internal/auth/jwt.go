package auth

import (
	"github.com/golang-jwt/jwt/v5"
	"github.com/pkg/errors"
	"os"
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
