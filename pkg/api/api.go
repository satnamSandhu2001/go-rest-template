package api

import (
	"encoding/json"
	"go-rest-template/pkg/config"
	"net/http"
	"time"
)

type Response struct {
	Success bool   `json:"success"`
	Message string `json:"message,omitempty"`
	Data    any    `json:"data,omitempty"`
	Errors  any    `json:"errors,omitempty"`
}

// sends a successful JSON response
func Success(w http.ResponseWriter, message string, data any) {
	writeJSON(w, http.StatusOK, Response{
		Success: true,
		Message: message,
		Data:    data,
	})
}

// sends an error JSON response with status 400
func Error(w http.ResponseWriter, message string) {
	writeJSON(w, http.StatusBadRequest, Response{
		Success: false,
		Message: message,
	})
}

// sends a 422 JSON response with validation error map
func ValidationErrors(w http.ResponseWriter, errors map[string]string) {
	writeJSON(w, http.StatusUnprocessableEntity, Response{
		Success: false,
		Errors:  errors,
	})
}

// sends a JSON response with a custom status code
func AbortWithStatusError(w http.ResponseWriter, code int, message string) {
	writeJSON(w, code, Response{
		Success: false,
		Message: message,
	})
}

// sends a 404 response
func NotFound(w http.ResponseWriter, message string) {
	AbortWithStatusError(w, http.StatusNotFound, message)
}

// logs the error and sends a 500 response
func InternalServerError(w http.ResponseWriter) {
	AbortWithStatusError(w, http.StatusInternalServerError, "Internal server error")
}

// sends a 401 response
func Unauthorized(w http.ResponseWriter, message string) {
	AbortWithStatusError(w, http.StatusUnauthorized, message)
}

// sends a 403 response
func Forbidden(w http.ResponseWriter, message string) {
	AbortWithStatusError(w, http.StatusForbidden, message)
}

// sets a secure cookie and sends a success response
func SendJWTtoken(w http.ResponseWriter, r *http.Request, token string, message string, data any) {
	cookie := &http.Cookie{
		Name:     "access_token",
		Value:    token,
		Path:     "/",
		Expires:  time.Now().Add(time.Duration(config.APP().COOKIE_AGE_HOURS) * time.Hour),
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
	}

	if config.APP().GO_ENV == "production" {
		cookie.Secure = true
		if config.APP().COOKIE_DOMAIN != "" {
			cookie.Domain = config.APP().COOKIE_DOMAIN
		}
	}

	http.SetCookie(w, cookie)
	Success(w, message, data)
}

// removes the JWT cookie
func Logout(w http.ResponseWriter) {
	cookie := &http.Cookie{
		Name:     "access_token",
		Value:    "",
		Path:     "/",
		MaxAge:   -1,
		HttpOnly: true,
		SameSite: http.SameSiteStrictMode,
	}
	http.SetCookie(w, cookie)
	Success(w, "Logged out successfully", nil)
}

// a helper function to write a JSON response
func writeJSON(w http.ResponseWriter, statusCode int, resp Response) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(resp)
}
