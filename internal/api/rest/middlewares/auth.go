package middlewares

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"booking-service/internal/api/rest/helpers"
	"github.com/golang-jwt/jwt/v5"
)

const (
	missingAuthHeader   = "Missing Authorization header"
	invalidHeaderFormat = "Invalid Authorization header format"
	invalidToken        = "Invalid or expired token"
	invalidTokenClaims  = "Invalid token claims"
	invalidTokenPayload = "Invalid token payload"
)

type JWTMiddleware struct {
	jwtSecret string
}

func NewJWTMiddleware(secret string) *JWTMiddleware {
	return &JWTMiddleware{jwtSecret: secret}
}

func (m *JWTMiddleware) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(resp http.ResponseWriter, req *http.Request) {
		authHeader := req.Header.Get("Authorization")
		if authHeader == "" {
			helpers.WriteErrorResponse(resp, helpers.NewErrorResponse(missingAuthHeader, helpers.MissingAuthHeaderErr), http.StatusUnauthorized)
			return
		}

		// Header should be: "Authorization": "Bearer ..."
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			helpers.WriteErrorResponse(resp, helpers.NewErrorResponse(invalidHeaderFormat, helpers.InvalidHeaderFormatErr), http.StatusUnauthorized)
			return
		}

		tokenString := parts[0]

		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			// Validate the algorithm
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			return m.jwtSecret, nil
		})

		if err != nil {
			helpers.WriteErrorResponse(resp, helpers.NewErrorResponse(invalidToken, helpers.InvalidTokenErr), http.StatusUnauthorized)
			return
		}

		// Extract claims and put userID into context
		if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
			userID, ok := claims["user_id"].(string)
			if !ok {
				helpers.WriteErrorResponse(resp, helpers.NewErrorResponse(invalidTokenPayload, helpers.InvalidTokenErr), http.StatusUnauthorized)
				return
			}

			ctx := context.WithValue(req.Context(), "userID", userID)
			next.ServeHTTP(resp, req.WithContext(ctx))
		} else {
			helpers.WriteErrorResponse(resp, helpers.NewErrorResponse(invalidTokenClaims, helpers.InvalidTokenErr), http.StatusUnauthorized)
		}
	})
}
