package controllers

import (
	"net/http"

	"github.com/0xMarvell/scissor/pkg/config"
	"github.com/0xMarvell/scissor/pkg/models"
	"github.com/gin-gonic/gin"
)

func SayHello(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"msg": "Hello and welcome to scissor! shorten urls with ease!",
	})
}

func Shorten(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"msg": "get original url and shorten it",
	})
}

func Redirect(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"msg": "get shortened url and redirect",
	})
}

func GetURLs(c *gin.Context) {
	// Retireve all user objects from database
	var urls []models.URL
	config.DB.Find(&urls)
	// Get count of users in database
	var count int64
	result := config.DB.Model(&models.URL{}).Count(&count)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to get count",
		})
		return
	}
	// Return user details as JSON response
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"count":   count,
		"urls":    urls,
	})
}
