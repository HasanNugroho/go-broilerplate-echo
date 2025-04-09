package middleware

import (
	"net/http"
	"slices"

	"github.com/HasanNugroho/starter-golang/internal/shared/utils"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func CheckAccess(permission []string) gin.HandlerFunc {
	return func(c *gin.Context) {
		claimsRaw, exists := c.Get("claims")
		if !exists {
			utils.SendError(c, http.StatusForbidden, "No claims found", nil)
			c.Abort()
			return
		}

		// Parsing claims
		claims, ok := claimsRaw.(jwt.MapClaims)
		if !ok {
			utils.SendError(c, http.StatusForbidden, "Invalid claims format", nil)
			c.Abort()
			return
		}

		data, ok := claims["data"].(map[string]interface{})
		if !ok {
			utils.SendError(c, http.StatusForbidden, "Invalid data in claims", nil)
			c.Abort()
			return
		}

		rawRoles, ok := data["permission"].([]interface{})
		if !ok {
			utils.SendError(c, http.StatusForbidden, "Roles not found or wrong format", nil)
			c.Abort()
			return
		}

		// Convert []interface{} to []string
		var roles []string
		for _, role := range rawRoles {
			if str, ok := role.(string); ok {
				roles = append(roles, str)
			}
		}

		// Admin full access
		if slices.Contains(roles, "manage:system") {
			c.Next()
			return
		}

		// Check permission
		if len(utils.Intersection(roles, permission)) < 1 {
			utils.SendError(c, http.StatusForbidden, "Access denied", nil)
			c.Abort()
			return
		}

		c.Next()
	}
}
