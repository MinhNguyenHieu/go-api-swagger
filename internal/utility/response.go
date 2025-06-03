package utility

import (
	"encoding/json"
	"net/http"
)

func JSONResponse(w http.ResponseWriter, statusCode int, payload interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	if payload != nil {
		json.NewEncoder(w).Encode(payload)
	}
}

func ErrorResponse(w http.ResponseWriter, statusCode int, message string) {
	JSONResponse(w, statusCode, map[string]string{"message": message})
}
