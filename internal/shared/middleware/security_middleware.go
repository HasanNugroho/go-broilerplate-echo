package middleware

import (
	"net/http"

	"github.com/HasanNugroho/starter-golang/config"
	"github.com/HasanNugroho/starter-golang/internal/shared/utils"
	"github.com/gin-gonic/gin"
)

func SecurityMiddleware(config *config.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.Request.Host != config.Security.ExpectedHost {
			utils.SendError(c, http.StatusBadRequest, "Invalid host header", nil)
			c.Abort()
			return
		}
		c.Header("X-Frame-Options", config.Security.XFrameOptions)
		c.Header("Content-Security-Policy", config.Security.ContentSecurity)
		c.Header("X-XSS-Protection", config.Security.XXSSProtection)
		c.Header("Strict-Transport-Security", config.Security.StrictTransport)
		c.Header("Referrer-Policy", config.Security.ReferrerPolicy)
		c.Header("X-Content-Type-Options", config.Security.XContentTypeOpts)
		c.Header("Permissions-Policy", config.Security.PermissionsPolicy)
		c.Next()
	}
}
