package context

import "github.com/gin-gonic/gin"

func abstractGetter[T any](ctx *gin.Context, key string) T {
	data := ctx.Keys[key].(T)
	return data
}

func GetBody[T any](ctx *gin.Context) T {
	return abstractGetter[T](ctx, "body")
}

func GetURI[T any](ctx *gin.Context) T {
	return abstractGetter[T](ctx, "uri")
}

func GetQuery[T any](ctx *gin.Context) T {
	return abstractGetter[T](ctx, "query")
}
