package middleware

import (
	"backend/utils"
	"net/http"

	"github.com/gin-gonic/gin"
)

func DBErrorHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		if len(c.Errors) == 0 {
			return
		}

		for _, e := range c.Errors {
			if e.Err == nil {
				continue
			}
			status, payload := utils.ParseDBError(e.Err)
			if status >= 400 {
				c.AbortWithStatusJSON(status, utils.ErrorResponse("error", payload))
				return
			}
		}

		// If here, there were errors but none matched DB parsing -> generic 500
		c.AbortWithStatusJSON(http.StatusInternalServerError, utils.ErrorResponse("error", "internal server error"))
	}
}
