package middleware

import (
	"fmt"
	"net/http"

	"external-backend-go/internal/logger"
	"external-backend-go/internal/store"
	"external-backend-go/internal/utility"
)

func AuthRoleMiddleware(requiredRole string, userStore store.UserStore, roleStore store.RoleStore, appLogger *logger.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			claims, ok := GetUserClaimsFromContext(r.Context())
			if !ok {
				utility.UnauthorizedErrorResponse(w, r, fmt.Errorf("User claims not found in context"), appLogger)
				return
			}

			userRoleFromClaims, hasRoleInClaims := claims["role"].(string)
			if hasRoleInClaims && userRoleFromClaims == requiredRole {
				next.ServeHTTP(w, r)
				return
			}

			userIDFloat, ok := claims["sub"].(float64)
			if !ok {
				utility.ForbiddenResponse(w, r, appLogger)
				return
			}
			userID := int32(userIDFloat)

			dbUser, err := userStore.GetUserByID(r.Context(), userID)
			if err != nil {
				appLogger.Error("Failed to get user from DB for role check: %v", err)
				utility.ForbiddenResponse(w, r, appLogger)
				return
			}

			role, err := roleStore.GetByID(r.Context(), dbUser.RoleID)
			if err != nil {
				appLogger.Error("Failed to get role name from DB for role ID %d: %v", dbUser.RoleID, err)
				utility.ForbiddenResponse(w, r, appLogger)
				return
			}

			if role.Name != requiredRole {
				utility.ForbiddenResponse(w, r, appLogger)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}
