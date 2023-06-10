package controllers

import (
	"fmt"
	"net/http"

	"github.com/0xMarvell/scissor/pkg/config"
	"github.com/0xMarvell/scissor/pkg/models"
	"github.com/0xMarvell/scissor/pkg/utils"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// SayHello displays a simple greeting on the index page of the API
func SayHello(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"msg": "Hello and welcome to scissor! shorten urls with ease!",
	})
}

// Shorten takes the original url and shortens it
func Shorten(c *gin.Context) {
	var urlPayload struct {
		URL string `json:"url"`
	}
	if err := c.ShouldBindJSON(&urlPayload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Bad request"})
		return
	}

	// // Validate the URL
	if !utils.IsValidURL(urlPayload.URL) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid URL"})
		return
	}

	//Get UserID of the authenticated user
	userID, _ := c.Get("userID")

	// Generate short link
	shortURL := utils.GenerateShortLink(urlPayload.URL, fmt.Sprintf("%v", userID))

	// Save URL object to database
	result := config.DB.Create(&models.URL{
		OriginalURL:  urlPayload.URL,
		ShortenedURL: shortURL,
		UserID:       utils.UserIDToUint(userID),
	})
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
		return
	}

	// Return short URL
	host := "http://localhost:8080/"
	c.JSON(http.StatusOK, gin.H{
		"message":   "short url created successfully",
		"short_url": host + shortURL,
	})
}

// Redirect takes the shortened URL and
// redirects the user the the original URL
func Redirect(c *gin.Context) {
	userID, _ := c.Get("userID")
	var url models.URL
	result := config.DB.Where("user_id = ?", userID).First(&url)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "URL not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
		}
		return
	}
	c.Redirect(http.StatusMovedPermanently, url.OriginalURL)
}

// GetURLs retrieves all existing urls
// for the logged in user
func GetURLs(c *gin.Context) {
	// Retireve all user objects from database
	userID, _ := c.Get("userID")
	var urls []models.URL
	config.DB.Where("user_id = ?", userID).Find(&urls)
	// Get count of users in database
	var count int64
	result := config.DB.Model(&models.URL{}).Where("user_id = ?", userID).Count(&count)
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

// DeleteURL deletes an existing URL object from the database
func DeleteURL(c *gin.Context) {
	urlID := c.Param("urlID")
	userID, _ := c.Get("userID")
	var url models.URL

	config.DB.Where("id = ? AND user_id = ?", urlID, userID).First(&url)

	if !urlExists(urlID) {
		c.IndentedJSON(http.StatusNotFound, gin.H{
			"success": false,
			"error":   fmt.Sprintf("URL object with id {%v} does not exist", urlID),
		})
		return
	}
	if err := config.DB.Unscoped().Delete(&url).Error; err != nil {
		// Handle deletion error
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to delete URL - you cannot delete a URL you did not create"})
		return
	}

	c.JSON(http.StatusNoContent, gin.H{})
}

// urlExists checks if requested blog post exists in the database.
func urlExists(id string) bool {
	var url models.URL

	config.DB.First(&url, id)
	if url.ID == 0 {
		return false
	} else {
		return true
	}
}