package lang

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/skatekrak/scribe/database"
)

func FindAll(c *gin.Context) {
	db, _ := database.Database()

	var langs []Lang
	if err := db.Find(&langs).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, langs)
}

type CreateBody struct {
	IsoCode  string `json:"isoCode" binding:"required,len=2"`
	ImageURL string `json:"imageURL" binding:"required"`
}

func Create(c *gin.Context) {
	var body CreateBody
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Wrong body",
			"error":   err.Error(),
		})
		return
	}

	log.Println(body)

	db, _ := database.Database()

	lang := Lang{
		IsoCode:  body.IsoCode,
		ImageURL: body.ImageURL,
	}

	if err := db.Create(&lang).Error; err != nil {
		c.JSON(http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusOK, lang)
}

type updateBody struct {
	ImageURL string `json:"imageURL" binding:"required"`
}

type langUri struct {
	IsoCode string `uri:"isoCode" binding:"required"`
}

func Update(c *gin.Context) {
	var uri langUri
	if err := c.ShouldBindUri(&uri); err != nil {
		c.JSON(http.StatusBadRequest, err)
	}

	log.Println("URI:", uri)

	var body updateBody
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Wrong body",
			"error":   err.Error(),
		})
		return
	}

	db, _ := database.Database()

	var lang Lang
	if err := db.Where("iso_code = ?", uri.IsoCode).First(&lang).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"message": "Lang not found",
		})
		return
	}

	lang.ImageURL = body.ImageURL
	if err := db.Save(&lang).Error; err != nil {
		c.JSON(http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusOK, lang)
}

func Delete(c *gin.Context) {
	var uri langUri
	if err := c.BindUri(&uri); err != nil {
		c.JSON(http.StatusBadRequest, err)
		return
	}

	db, _ := database.Database()

	var lang Lang
	if err := db.Where("iso_code = ?", uri.IsoCode).First(&lang).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"message": "Lang not found",
		})
		return
	}

	if err := db.Delete(&lang).Error; err != nil {
		c.JSON(http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Lang deleted",
	})
}
