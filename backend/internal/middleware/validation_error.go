package middleware

import (
	"backend/utils"
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

func ValidationErrorHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		fmt.Println(c.Request.Body)
		c.Next()
		if len(c.Errors) > 0 {
			for _, e := range c.Errors {
				if e.Type == gin.ErrorTypeBind {
					log.Printf("Validation error: %v\n", e.Err)
					errorsMap := utils.ValidationErrorHandler(e.Err)
					c.AbortWithStatusJSON(http.StatusBadRequest, utils.ErrorResponse("validation error", errorsMap))
					return
				}
			}
		}
	}
}
