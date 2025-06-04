package API

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

// Response represents a standard API Response structure for gin
type Response struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
	Data    any    `json:"data,omitempty"`
}

// Success sends a successful JSON response
func Success(c *gin.Context, message string, data any) {
	c.JSON(http.StatusOK, Response{
		Success: true,
		Message: message,
		Data:    data,
	})
}

// Error sends an error JSON Response with status 400 and message
func Error(c *gin.Context, message string) {
	c.JSON(http.StatusBadRequest, Response{
		Success: false,
		Message: message,
	})
}

// ValidationsErrors returns an error JSON Response
func ValidationsErrors(c *gin.Context, errors map[string]string) {
	c.JSON(http.StatusUnprocessableEntity, Response{
		Success: false,
		Data:    errors,
	})
}

// AbortWithStatusError sends an error JSON response with status and aborts the request
func AbortWithStatusError(c *gin.Context, code int, message string) {
	c.AbortWithStatusJSON(code, Response{
		Success: false,
		Message: message,
	})
}

// NotFound sends a 404 Not Found response and aborts
func NotFound(c *gin.Context, message string) {
	AbortWithStatusError(c, http.StatusNotFound, message)
}

// InternalServerError logs error and sends a 500 Internal Server Error response and aborts
func InternalServerError(c *gin.Context, message string, logs error) {
	log.Println("Internal Server Error :", logs)
	AbortWithStatusError(c, http.StatusInternalServerError, message)
}

// Unauthorized sends a 401 Unauthorized response and aborts
func Unauthorized(c *gin.Context, message string) {
	AbortWithStatusError(c, http.StatusUnauthorized, message)
}

// Forbidden sends a 403 Forbidden response and aborts
func Forbidden(c *gin.Context, message string) {
	AbortWithStatusError(c, http.StatusForbidden, message)
}

// SendJWTtoken sends a JWT token in a cookie with JSON Response
func SendJWTtoken(c *gin.Context, token string, message string, data any) {
	c.SetSameSite(http.SameSiteLaxMode)
	c.SetCookie("Authorization", "Bearer "+token, 3400*24*30, "", "", false, true)
	Success(c, message, data)
}
