package middleware

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt/v5"

	"external-backend-go/internal/auth"
	"external-backend-go/internal/logger"
	"external-backend-go/internal/utility"
)

type contextKey string

const userClaimsContextKey contextKey = "userClaims"

func AuthMiddleware(jwtSecret string, appLogger *logger.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			tokenString := r.Header.Get("Authorization")
			if tokenString == "" {
				utility.UnauthorizedErrorResponse(w, r, fmt.Errorf("Authentication token required"), appLogger)
				return
			}

			if !strings.HasPrefix(tokenString, "Bearer ") {
				utility.UnauthorizedErrorResponse(w, r, fmt.Errorf("Invalid token format"), appLogger)
				return
			}
			tokenString = strings.TrimPrefix(tokenString, "Bearer ")

			claims, err := auth.ValidateToken(tokenString, jwtSecret)
			if err != nil {
				utility.UnauthorizedErrorResponse(w, r, fmt.Errorf("Invalid token: %w", err), appLogger)
				return
			}

			ctx := context.WithValue(r.Context(), userClaimsContextKey, claims)
			r = r.WithContext(ctx)

			next.ServeHTTP(w, r)
		})
	}
}

func GetUserClaimsFromContext(ctx context.Context) (jwt.MapClaims, bool) {
	claims, ok := ctx.Value(userClaimsContextKey).(jwt.MapClaims)
	return claims, ok
}
