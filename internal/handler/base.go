package handler

import (
	"net/http"

	"external-backend-go/internal/logger"
	"external-backend-go/internal/utility"
)

type BaseHandler struct {
	Logger *logger.Logger
}

func NewBaseHandler(logger *logger.Logger) *BaseHandler {
	return &BaseHandler{Logger: logger}
}

func (b *BaseHandler) JSONResponse(w http.ResponseWriter, statusCode int, payload interface{}) {
	utility.JSONResponse(w, statusCode, payload)
}

func (b *BaseHandler) ErrorResponse(w http.ResponseWriter, statusCode int, message string) {
	utility.ErrorResponse(w, statusCode, message)
}

// func (b *BaseHandler) InternalServerError(w http.ResponseWriter, r *http.Request, err error) {
// 	utility.InternalServerError(w, r, err)
// }

// func (b *BaseHandler) BadRequestResponse(w http.ResponseWriter, r *http.Request, err error) {
// 	utility.BadRequestResponse(w, r, err)
// }

// func (b *BaseHandler) NotFoundResponse(w http.ResponseWriter, r *http.Request) {
// 	utility.NotFoundResponse(w, r)
// }

// func (b *BaseHandler) UnauthorizedErrorResponse(w http.ResponseWriter, r *http.Request, err error) {
// 	utility.UnauthorizedErrorResponse(w, r, err)
// }

// func (b *BaseHandler) ForbiddenResponse(w http.ResponseWriter, r *http.Request) {
// 	utility.ForbiddenResponse(w, r)
// }
