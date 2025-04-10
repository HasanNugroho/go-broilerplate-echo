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
			// Recover from panic
			defer func() {
				if err := recover(); err != nil {
					app.Log.Error().Msgf("Recovered from panic: %v", err)
					utils.SendError(c, http.StatusInternalServerError, "Internal Server Error", nil)
				}
			}()

			err := next(c)
			if err == nil {
				return nil
			}

			// Handle error custom
			switch e := err.(type) {
			case *utils.BadRequestError:
				utils.SendError(c, http.StatusBadRequest, e.Message, nil)
			case *utils.UnauthorizedError:
				utils.SendError(c, http.StatusUnauthorized, e.Message, nil)
			case *utils.ForbiddenError:
				utils.SendError(c, http.StatusForbidden, e.Message, nil)
			case *utils.NotFoundError:
				utils.SendError(c, http.StatusNotFound, e.Message, nil)
			case *utils.ConflictError:
				utils.SendError(c, http.StatusConflict, e.Message, nil)
			case *utils.InternalError:
				utils.SendError(c, http.StatusInternalServerError, e.Message, nil)
			default:
				app.Log.Error().Err(err).Msg("Unhandled error")
				utils.SendError(c, http.StatusInternalServerError, "Internal Server Error", nil)
			}

			return nil
		}
	}
}
