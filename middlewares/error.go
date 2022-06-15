package middlewares

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/skatekrak/scribe/formatter"
)

func ErrorHandler() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ctx.Next()

		for _, ginErr := range ctx.Errors {
			if verr, ok := ginErr.Err.(validator.ValidationErrors); ok {
				f := formatter.NewJSONFormatter()
				ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
					"error": f.Simple(verr),
				})
				return
			}
		}
	}
}
