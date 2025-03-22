package response

import (
	"github.com/gin-gonic/gin"
)

// Response format untuk sukses dan error
type Response struct {
	Status  int         `json:"status"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
	Error   interface{} `json:"error,omitempty"`
}

// SendSuccess mengirim response sukses
func SendSuccess(c *gin.Context, statusCode int, message string, data interface{}) {
	c.JSON(statusCode, Response{
		Status:  statusCode,
		Message: message,
		Data:    data,
	})
}

// SendError mengirim response error
func SendError(c *gin.Context, statusCode int, message string, err interface{}) {
	c.JSON(statusCode, Response{
		Status:  statusCode,
		Message: message,
		Error:   err,
	})
}
