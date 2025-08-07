package middleware

import (
	"context"
	"fmt"
	"strings"

	"github.com/go-kratos/kratos/v2/errors"
	"github.com/go-kratos/kratos/v2/middleware"
	"github.com/go-kratos/kratos/v2/transport"
	"github.com/golang-jwt/jwt/v4"
)

var (
	// A hardcoded secret for demonstration.
	// In a real production environment, this should be loaded from a secure configuration source.
	jwtSecret = []byte("a_very_secret_key_that_should_be_long_and_random")
)

// JWTAuthMiddleware creates a JWT authentication middleware.
func JWTAuthMiddleware() middleware.Middleware {
	return func(handler middleware.Handler) middleware.Handler {
		return func(ctx context.Context, req interface{}) (reply interface{}, err error) {
			// Attempt to retrieve transport context
			tr, ok := transport.FromServerContext(ctx)
			if !ok {
				// If it's not a server transport context, we can't get headers.
				// Depending on the policy, we might block or allow the request.
				// For an API gateway, it's safer to block.
				return nil, errors.Unauthorized("UNAUTHORIZED", "Could not retrieve transport context")
			}

			// Get the Authorization header from the request.
			authHeader := tr.RequestHeader().Get("Authorization")
			if authHeader == "" {
				return nil, errors.Unauthorized("UNAUTHORIZED", "Authorization header is required")
			}

			// The header should be in the format "Bearer {token}"
			parts := strings.Split(authHeader, " ")
			if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") {
				return nil, errors.Unauthorized("UNAUTHORIZED", "Authorization header format must be 'Bearer {token}'")
			}
			tokenString := parts[1]

			// Parse and validate the JWT token.
			token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
				// Ensure the signing method is what we expect.
				if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
					return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
				}
				return jwtSecret, nil
			})

			if err != nil {
				return nil, errors.Unauthorized("UNAUTHORIZED", "Invalid token: "+err.Error())
			}

			if _, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
				// The token is valid. Proceed with the request.
				// In a real app, you might add the user's claims to the context here.
				return handler(ctx, req)
			}

			return nil, errors.Unauthorized("UNAUTHORIZED", "Invalid token claims")
		}
	}
}
