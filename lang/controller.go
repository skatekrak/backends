package lang

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type Controller struct {
	db *gorm.DB
}

func NewController(db *gorm.DB) *Controller {
	return &Controller{db}
}

func (c *Controller) FindAll(ctx *gin.Context) {
	var langs []Lang
	if err := c.db.Find(&langs).Error; err != nil {
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
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": "Wrong body",
			"error":   err.Error(),
		})
		return
	}

	log.Println(body)

	lang := Lang{
		IsoCode:  body.IsoCode,
		ImageURL: body.ImageURL,
	}

	if err := c.db.Create(&lang).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, err)
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
		ctx.JSON(http.StatusBadRequest, err)
		return
	}

	log.Println("URI:", uri)

	var body updateBody
	if err := ctx.ShouldBindJSON(&body); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": "Wrong body",
			"error":   err.Error(),
		})
		return
	}

	var lang Lang
	if err := c.db.Where("iso_code = ?", uri.IsoCode).First(&lang).Error; err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{
			"message": "Lang not found",
		})
		return
	}

	lang.ImageURL = body.ImageURL
	if err := c.db.Save(&lang).Error; err != nil {
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

	var lang Lang
	if err := c.db.Where("iso_code = ?", uri.IsoCode).First(&lang).Error; err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{
			"message": "Lang not found",
		})
		return
	}

	if err := c.db.Delete(&lang).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, err)
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "Lang deleted",
	})
}
