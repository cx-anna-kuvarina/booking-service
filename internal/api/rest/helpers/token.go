package helpers

import (
	"encoding/base64"
	"encoding/json"
	"strings"

	"github.com/pkg/errors"
)

const (
	jwtLength = 3
)

type Claims struct {
	UserID string `json: "user_id"`
}

func ExtractClaims(token string) (*Claims, error) {
	tokenSlice := strings.Split(token, ".")
	if len(tokenSlice) < jwtLength {
		return nil, errors.New("token is not valid, missing parts")
	}

	decoded, err := base64.RawURLEncoding.DecodeString(tokenSlice[1])
	if err != nil {
		return nil, errors.Wrap(err, "could not decode token")
	}

	var claims Claims
	if err := json.Unmarshal(decoded, &claims); err != nil {
		return nil, errors.Wrap(err, "could not unmarshal token")
	}

	return &claims, nil
}

func GetUserIDFromJWT(token string) (string, error) {
	claims, err := ExtractClaims(token)
	if err != nil {
		return "", errors.Wrap(err, "could not extract claims from token")
	}
	return claims.UserID, nil
}
