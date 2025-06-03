package middleware

import (
	"net/http"
)

func BasicAuthMiddleware(username, password string, errorHandler func(w http.ResponseWriter, r *http.Request, err error)) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			user, pass, ok := r.BasicAuth()

			if !ok || !authenticateBasic(user, pass, username, password) {
				errorHandler(w, r, nil)
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}

func authenticateBasic(user, pass, expectedUser, expectedPass string) bool {
	return user == expectedUser && pass == expectedPass
}
