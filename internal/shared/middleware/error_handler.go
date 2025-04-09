package middleware

import (
	"net/http"

	"github.com/HasanNugroho/starter-golang/internal/app"
	"github.com/HasanNugroho/starter-golang/internal/shared/utils"
	"github.com/labstack/echo/v4"
)

func ErrorHandler(app *app.Apps) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			defer func() {
				if err := recover(); err != nil {
					app.Log.Error().Msgf("Recovered from panic: %v", err)
					utils.SendError(c, http.StatusInternalServerError, "Internal Server Error", nil)
				}
			}()

			return next(c)
		}
	}
}
