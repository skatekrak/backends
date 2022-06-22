package lang

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/skatekrak/scribe/context"
	"github.com/skatekrak/scribe/model"
	"gorm.io/gorm"
)

type Controller struct {
	s *Service
}

func NewController(db *gorm.DB) *Controller {
	return &Controller{s: NewService(db)}
}

func (c *Controller) LoaderHandler() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		uri := context.GetURI[LangUri](ctx)

		lang, err := c.s.Get(uri.IsoCode)
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusNotFound, gin.H{
				"message": "Lang not found",
			})
			return
		}

		ctx.Set("lang", lang)
		ctx.Next()
	}
}

func (c *Controller) FindAll(ctx *gin.Context) {
	langs, err := c.s.FindAll()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, langs)
}

func (c *Controller) Create(ctx *gin.Context) {
	body := context.GetBody[CreateBody](ctx)

	if _, err := c.s.Get(body.IsoCode); err == nil {
		ctx.AbortWithStatusJSON(http.StatusConflict, gin.H{
			"message": "isoCode already used",
		})
		return
	}

	lang := model.Lang{
		IsoCode:  body.IsoCode,
		ImageURL: body.ImageURL,
	}

	if err := c.s.Create(&lang); err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, err)
		return
	}

	ctx.JSON(http.StatusOK, lang)
}

func (c *Controller) Update(ctx *gin.Context) {
	lang := ctx.Keys["lang"].(model.Lang)
	body := context.GetBody[UpdateBody](ctx)

	lang.ImageURL = body.ImageURL
	if err := c.s.Update(&lang); err != nil {
		ctx.JSON(http.StatusInternalServerError, err)
		return
	}

	ctx.JSON(http.StatusOK, lang)
}

func (c *Controller) Delete(ctx *gin.Context) {
	lang := ctx.Keys["lang"].(model.Lang)

	if err := c.s.Delete(&lang); err != nil {
		ctx.JSON(http.StatusInternalServerError, err)
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "Lang deleted",
	})
}
