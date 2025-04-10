package middleware

import (
	"net/http"

	"github.com/HasanNugroho/starter-golang/internal/app"
	"github.com/HasanNugroho/starter-golang/internal/shared/utils"
	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
)

func AuthMiddleware(app *app.Apps) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			tokenString := c.Request().Header.Get("Authorization")

			token, err := utils.ValidateToken(app, tokenString)
			if err != nil || !token.Valid {
				utils.SendError(c, http.StatusUnauthorized, "Unauthorized", nil)
				return nil

			}

			if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
				c.Set("claims", claims)
			} else {
				utils.SendError(c, http.StatusUnauthorized, "Unauthorized", nil)
				return nil
			}

			return next(c)
		}
	}
}
