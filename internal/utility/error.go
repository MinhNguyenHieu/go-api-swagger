package utility

import (
	"fmt"
	"net/http"

	"external-backend-go/internal/logger"
)

func UnauthorizedErrorResponse(w http.ResponseWriter, r *http.Request, err error, appLogger *logger.Logger) {
	appLogger.Warn("Unauthorized access attempt: %v", err)
	ErrorResponse(w, http.StatusUnauthorized, fmt.Sprintf("Authentication failed: %v", err))
}

func UnauthorizedBasicErrorResponse(w http.ResponseWriter, r *http.Request, err error, appLogger *logger.Logger) {
	w.Header().Set("WWW-Authenticate", "Basic realm=\"restricted\"")
	appLogger.Warn("Basic Auth failed: %v", err)
	ErrorResponse(w, http.StatusUnauthorized, fmt.Sprintf("Basic authentication failed: %v", err))
}

func InternalServerError(w http.ResponseWriter, r *http.Request, err error, appLogger *logger.Logger) {
	appLogger.Error("Internal server error: %v", err)
	ErrorResponse(w, http.StatusInternalServerError, "An internal server error occurred.")
}

func ForbiddenResponse(w http.ResponseWriter, r *http.Request, appLogger *logger.Logger) {
	appLogger.Warn("Forbidden access attempt for %s %s", r.Method, r.URL.Path)
	ErrorResponse(w, http.StatusForbidden, "You do not have permission to access this resource.")
}

func RateLimitExceededResponse(w http.ResponseWriter, r *http.Request, retryAfter string, appLogger *logger.Logger) {
	w.Header().Set("Retry-After", retryAfter)
	appLogger.Warn("Rate limit exceeded for %s", r.RemoteAddr)
	ErrorResponse(w, http.StatusTooManyRequests, "You have exceeded the request limit. Please try again later.")
}

func BadRequestResponse(w http.ResponseWriter, r *http.Request, err error, appLogger *logger.Logger) {
	appLogger.Warn("Bad request: %v", err)
	ErrorResponse(w, http.StatusBadRequest, fmt.Sprintf("Invalid request: %v", err))
}

func NotFoundResponse(w http.ResponseWriter, r *http.Request, appLogger *logger.Logger) {
	appLogger.Warn("Resource not found: %s %s", r.Method, r.URL.Path)
	ErrorResponse(w, http.StatusNotFound, "Resource not found.")
}
