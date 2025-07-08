package token

import (
	"booking-service/internal/api/rest/helpers"
	"strings"
)

// ExtractUserIDFromJWT extracts the user ID from a JWT token string (with or without Bearer prefix)
func ExtractUserIDFromJWT(token string) (string, error) {
	if len(token) > 7 && strings.HasPrefix(token, "Bearer ") {
		token = token[7:]
	}
	return helpers.GetUserIDFromJWT(token)
}
