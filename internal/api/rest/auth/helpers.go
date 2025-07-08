package auth

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func generateJWT(userId, jwtSecret string, expirationPeriod time.Duration) (string, error) {
	claims := jwt.MapClaims{
		"user_id": userId,
		"exp":     time.Now().Add(expirationPeriod).Unix(),
		"iat":     time.Now().Unix(),
	}
	// Create token with claims
	token := jwt.NewWithClaims(jwt.SigningMethodES512, claims)
	// Sign the token with jwt secret
	return token.SignedString(jwtSecret)
}
