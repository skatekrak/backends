package source

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/skatekrak/scribe/lang"
	"github.com/skatekrak/scribe/model"
	"gorm.io/gorm"
)

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

type findAllQuery struct {
	Types []string `form:"types[]" binding:"dive,eq=vimeo|eq=youtube|eq=rss"`
}

func (c *Controller) FindAll(ctx *gin.Context) {
	var query findAllQuery
	if err := ctx.BindQuery(&query); err != nil {
		return
	}

	log.Println(query)

	sources, err := c.s.FindAll()
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
	}

	ctx.JSON(http.StatusOK, sources)
}

type createBody struct {
	Type        string `json:"type" binding:"required"`
	LangID      string `json:"langID" binding:"required"`
	SkateSource bool   `json:"skateSource" default:"false"`
	SourceID    string `json:"sourceID" binding:"required"`
}

func (c *Controller) Create(ctx *gin.Context) {
	var body createBody
	if err := ctx.ShouldBindJSON(&body); err != nil {
		ctx.Error(err)
		return
	}

	if _, err := c.s.GetBySourceID(body.SourceID); err == nil {
		ctx.AbortWithStatusJSON(http.StatusConflict, gin.H{
			"message": "Source with this sourceID already in use",
		})
		return
	}

	nextOrder, err := c.s.GetNextOrder()
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"message": "Oops, couldn't compute the next order",
			"error":   err.Error(),
		})
		return
	}

	if _, errLang := c.ls.Get(body.LangID); errLang != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"message": "This langID doesn't exist",
		})
		return
	}

	source := model.Source{
		Order:       nextOrder,
		SourceType:  body.Type,
		SkateSource: body.SkateSource,
		SourceID:    body.SourceID,
		LangIsoCode: body.LangID,
	}

	if err := c.s.Create(&source); err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"message": "Something went wrong",
			"error":   err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, source)
}
