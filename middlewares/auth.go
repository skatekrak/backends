package middlewares

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type authHeader struct {
	Authorization string `header:"Authorization" binding:"required"`
}

func Authorization(apiKey string) gin.HandlerFunc {
	return func(c *gin.Context) {
		var header authHeader
		if err := c.BindHeader(&header); err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"message": "Missing Authorization header",
			})
			return
		}

		if header.Authorization != apiKey {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
				"message": "Wait, that's illegal",
			})
			return
		}

		c.Next()
	}
}
