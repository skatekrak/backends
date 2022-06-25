package source

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/skatekrak/scribe/lang"
	"github.com/skatekrak/scribe/middlewares"
	"github.com/skatekrak/scribe/model"
	"gorm.io/gorm"
)

// Key used to pass the source interface between middlewares
const context_source = "source"

type Controller struct {
	s  *Service
	ls *lang.Service
}

func NewController(db *gorm.DB) *Controller {
	return &Controller{
		s:  NewService(db),
		ls: lang.NewService(db),
	}
}

func (c *Controller) LoaderHandler() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		uri := ctx.Keys[middlewares.URI].(SourceURI)

		source, err := c.s.Get(uri.SourceID)
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusNotFound, gin.H{
				"message": "Source not found",
			})
			return
		}

		ctx.Set(context_source, source)
		ctx.Next()
	}
}

func (c *Controller) FindAll(ctx *gin.Context) {
	query := ctx.Keys[middlewares.QUERY].(FindAllQuery)

	sources, err := c.s.FindAll(query.Types)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
	}

	ctx.JSON(http.StatusOK, sources)
}

func (c *Controller) Create(ctx *gin.Context) {
	body := ctx.Keys[middlewares.BODY].(CreateBody)

	if _, err := c.s.GetBySourceID(body.ChannelID); err == nil {
		ctx.AbortWithStatusJSON(http.StatusConflict, gin.H{
			"message": "This youtube channel is already added",
		})
		return
	}

	nextOrder, err := c.s.GetNextOrder()
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"message": "Couldn't process the next order",
			"error":   err.Error(),
		})
		return
	}

	// TODO: fetch youtube info

	source := model.Source{
		Order:       nextOrder,
		SourceType:  body.Type,
		SkateSource: body.IsSkateSource,
		LangIsoCode: body.LangIsoCode,
		SourceID:    body.ChannelID,
	}

	if err := c.s.Create(&source); err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"message": "Couldn't create the source",
			"error":   err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, source)
}

func (c *Controller) Update(ctx *gin.Context) {
	body := ctx.Keys[middlewares.BODY].(UpdateBody)
	source := ctx.Keys[context_source].(model.Source)

	log.Println(body)

	ctx.JSON(http.StatusOK, gin.H{
		"body":   body,
		"source": source,
	})
}

func (c *Controller) Delete(ctx *gin.Context) {
	source := ctx.Keys[context_source].(model.Source)

	if err := c.s.Delete(&source); err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"message": "Couldn't delete this source",
			"error":   err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "Source deleted",
	})
}
