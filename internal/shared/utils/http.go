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

// func SendPagination(c echo.Context, status int, message string, items interface{}, page, limit, totalItems int, err interface{}) {
// 	if limit < 1 {
// 		limit = 10
// 	}
// 	if page < 1 {
// 		page = 1
// 	}

// 	totalPages := int(math.Ceil(float64(totalItems) / float64(limit)))

// 	// Struktur response dengan pagination
// 	response := shared.Response{
// 		Status:  status,
// 		Message: message,
// 		Data: shared.DataWithPagination{
// 			Items: items,
// 			Paging: shared.Pagination{
// 				Limit:      limit,
// 				Page:       page,
// 				TotalItems: totalItems,
// 				TotalPages: totalPages,
// 			},
// 		},
// 	}

// }
