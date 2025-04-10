package utils

import (
	shared "github.com/HasanNugroho/starter-golang/internal/shared/model"
	"github.com/labstack/echo/v4"
)

// SendSuccess mengirim response sukses
func SendSuccess(c echo.Context, statusCode int, message string, data interface{}) {
	c.JSON(statusCode, shared.Response{
		Status:  statusCode,
		Message: message,
		Data:    data,
	})
}

// SendError mengirim response error
func SendError(c echo.Context, statusCode int, message string, err interface{}) {
	c.JSON(statusCode, shared.Response{
		Status:  statusCode,
		Message: message,
		Data:    err,
	})
}
