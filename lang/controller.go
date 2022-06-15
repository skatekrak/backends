package lang

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type Controller struct {
	s *Service
}

func NewController(db *gorm.DB) *Controller {
	return &Controller{s: &Service{db}}
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

type createBody struct {
	IsoCode  string `json:"isoCode" binding:"required,len=2"`
	ImageURL string `json:"imageURL" binding:"required"`
}

func (c *Controller) Create(ctx *gin.Context) {
	var body createBody
	if err := ctx.ShouldBindJSON(&body); err != nil {
		ctx.Error(err)
		return
	}

	if _, err := c.s.Get(body.IsoCode); err == nil {
		ctx.AbortWithStatusJSON(http.StatusConflict, gin.H{
			"message": "isoCode already used",
		})
		return
	}

	lang := Lang{
		IsoCode:  body.IsoCode,
		ImageURL: body.ImageURL,
	}

	if err := c.s.Create(&lang); err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, err)
		return
	}

	ctx.JSON(http.StatusOK, lang)
}

type updateBody struct {
	ImageURL string `json:"imageURL" binding:"required"`
}

type langUri struct {
	IsoCode string `uri:"isoCode" binding:"required"`
}

func (c *Controller) Update(ctx *gin.Context) {
	var uri langUri
	if err := ctx.ShouldBindUri(&uri); err != nil {
		ctx.Error(err)
		return
	}

	var body updateBody
	if err := ctx.ShouldBindJSON(&body); err != nil {
		ctx.Error(err)
		return
	}

	lang, err := c.s.Get(uri.IsoCode)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{
			"message": "Lang not found",
		})
		return
	}

	lang.ImageURL = body.ImageURL
	if err := c.s.Update(&lang); err != nil {
		ctx.JSON(http.StatusInternalServerError, err)
		return
	}

	ctx.JSON(http.StatusOK, lang)
}

func (c *Controller) Delete(ctx *gin.Context) {
	var uri langUri
	if err := ctx.BindUri(&uri); err != nil {
		ctx.JSON(http.StatusBadRequest, err)
		return
	}

	lang, err := c.s.Get(uri.IsoCode)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{
			"message": "Lang not found",
		})
		return
	}

	if err := c.s.Delete(&lang); err != nil {
		ctx.JSON(http.StatusInternalServerError, err)
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "Lang deleted",
	})
}
