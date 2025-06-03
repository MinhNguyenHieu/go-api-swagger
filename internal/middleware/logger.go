package middleware

import (
	"net/http"
	"time"

	"external-backend-go/internal/logger"
)

func LoggerMiddleware(next http.Handler) http.Handler {
	appLogger := logger.NewLogger()
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		next.ServeHTTP(w, r)
		appLogger.Info("[%s] %s %s %s", r.Method, r.RequestURI, time.Since(start), r.RemoteAddr)
	})
}
