package utils

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// SuccessResponse defines the standard response structure for successful API calls.
type SuccessResponse struct {
	Status  int         `json:"status"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

// ErrorResponse defines the standard response structure for errors.
type ErrorResponse struct {
	Status  int         `json:"status"`
	Message string      `json:"message"`
	Error   interface{} `json:"error,omitempty"`
}

// SendSuccess sends a standardized JSON response for successful operations.
func SendSuccess(c *gin.Context, message string, data interface{}) {
	c.JSON(http.StatusOK, SuccessResponse{
		Status:  http.StatusOK,
		Message: message,
		Data:    data,
	})
}

// SendError sends a standardized JSON response for errors.
func SendError(c *gin.Context, statusCode int, message string, error interface{}) {
	c.JSON(statusCode, ErrorResponse{
		Status:  statusCode,
		Message: message,
		Error:   error,
	})
}
