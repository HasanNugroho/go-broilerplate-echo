package middleware

import (
	"net/http"

	"github.com/HasanNugroho/starter-golang/internal/app"
	"github.com/HasanNugroho/starter-golang/internal/shared/utils"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func AuthMiddleware(app *app.Apps) gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString := c.GetHeader("Authorization")

		// Parse the token
		token, err := utils.ValidateToken(app, tokenString)

		if err != nil || !token.Valid {
			utils.SendError(c, http.StatusUnauthorized, "Unauthorized", nil)
			c.Abort()
			return
		}

		// Set the token claims to the context
		if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
			c.Set("claims", claims)
		} else {
			utils.SendError(c, http.StatusUnauthorized, "Unauthorized", nil)
			c.Abort()
			return
		}

		c.Next() // Proceed to the next handler if authorized
	}
}
