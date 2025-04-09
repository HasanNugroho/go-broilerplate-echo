package middleware

import (
	"net/http"
	"slices"

	"github.com/HasanNugroho/starter-golang/internal/shared/utils"
	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
)

func CheckAccess(permission []string) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			claimsRaw := c.Get("claims")
			if claimsRaw == nil {
				utils.SendError(c, http.StatusForbidden, "No claims found", nil)
				return nil
			}

			claims, ok := claimsRaw.(jwt.MapClaims)
			if !ok {
				utils.SendError(c, http.StatusForbidden, "Invalid claims format", nil)
				return nil
			}

			data, ok := claims["data"].(map[string]interface{})
			if !ok {
				utils.SendError(c, http.StatusForbidden, "Invalid data in claims", nil)
				return nil
			}

			rawRoles, ok := data["permission"].([]interface{})
			// app.GlobalApps.Log.Info().Msgf("Roles: %v", rawRoles)
			if !ok {
				utils.SendError(c, http.StatusForbidden, "Roles not found or wrong format", nil)
				return nil
			}

			var roles []string
			for _, role := range rawRoles {
				if str, ok := role.(string); ok {
					roles = append(roles, str)
				}
			}

			// Allow if manage:system exists
			if slices.Contains(roles, "manage:system") {
				return next(c)
			}

			// Otherwise, check if any permission matches
			if len(utils.Intersection(roles, permission)) < 1 {
				utils.SendError(c, http.StatusForbidden, "Access denied", nil)
				return nil
			}

			return next(c)
		}
	}
}
