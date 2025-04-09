package middleware

import (
	"net/http"

	"github.com/HasanNugroho/starter-golang/internal/app"
	"github.com/HasanNugroho/starter-golang/internal/shared/utils"
	"github.com/gin-gonic/gin"
)

func ErrorHandler(app *app.Apps) gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				app.Log.Error().Msgf("Recovered from panic: %v", err)

				utils.SendError(c, http.StatusInternalServerError, "Internal Server Error", nil)
				c.Abort()
				return
			}
		}()

		c.Next()
	}
}
