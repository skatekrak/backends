package middlewares

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/skatekrak/scribe/formatter"
)

func JSONHandler[T any]() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var body T
		if err := ctx.ShouldBindJSON(&body); err != nil {
			if verr, ok := err.(validator.ValidationErrors); ok {
				f := formatter.NewJSONFormatter()
				ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
					"error": f.Simple(verr),
				})
				return
			}
			ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
				"message": "Couldn't parse or format error",
				"error":   err.Error(),
			})
			return
		}

		ctx.Set("body", body)
		ctx.Next()
	}
}

func URIHandler[T any]() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var uri T
		if err := ctx.BindUri(&uri); err != nil {
			if verr, ok := err.(validator.ValidationErrors); ok {
				f := formatter.NewJSONFormatter()
				ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
					"error": f.Simple(verr),
				})
				return
			}
			ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
				"message": "Couldn't parse or format error",
				"error":   err.Error(),
			})
			return
		}

		ctx.Set("uri", uri)
		ctx.Next()
	}
}
