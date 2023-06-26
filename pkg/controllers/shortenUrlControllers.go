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

func SayHello(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "Hello and welcome to Scissor! shorten URLs in the blink of an eye!",
	})
}

// Shorten takes the original url and shortens it
//
// @Description Get original url and created a shortened version
// @Tags urls
// @Accept json
// @Produce json
// @Param payload body models.ShortenURLRequest true "payload"
// @Success 200 {object} models.ShortenResponse
// @Failure 400 {object} models.ErrorResponse
// @Failue 500 {object} models.ErrorResponse
// @Router /shortener/shorten [post]
func Shorten(c *gin.Context) {
	var u models.ShortenURLRequest

	if err := c.ShouldBindJSON(&u); err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: "Bad request"})
		return
	}

	// // Validate the URL
	if !utils.IsValidURL(u.URL) {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: "Invalid URL"})
		return
	}

	//Get UserID of the authenticated user
	userID, _ := c.Get("userID")

	// Generate short link
	shortURL := utils.GenerateShortLink(u.URL, fmt.Sprintf("%v", userID))

	// Save URL object to database
	result := config.DB.Create(&models.URL{
		OriginalURL:  u.URL,
		ShortenedURL: shortURL,
		UserID:       utils.UserIDToUint(userID),
	})
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: "Database error"})
		return
	}

	// Return short URL
	// hostLocal := "http://localhost:8080/api/v1/" // for development
	host := "https://sci-ssor.onrender.com/"
	c.JSON(http.StatusOK, models.ShortenResponse{
		Message: "short url created successfully",
		// ShortURL: hostLocal + shortURL, // for development
		ShortURL: host + shortURL,
	})
}

// Redirect takes the shortened URL and
// redirects the user the the original URL
//
// @Description Redirect short url to original url
// @Tags urls
// @Produce json
// @Param shortURL path string true "Short URL"
// @Success 301 "Moved Permanently"
// @Failure 404 {object} models.ErrorResponse
// @Failue 500 {object} models.ErrorResponse
// @Router /shortener/redirect/{shortURL} [get]
func Redirect(c *gin.Context) {
	userID, _ := c.Get("userID")
	var url models.URL
	result := config.DB.Where("user_id = ?", userID).First(&url)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, models.ErrorResponse{Error: "URL not found"})
		} else {
			c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: "Database error"})
		}
		return
	}
	c.Redirect(http.StatusMovedPermanently, url.OriginalURL)
}

// GetURLs retrieves all existing urls
// for the logged in user
//
// @Description Retrieve all shortened urls created by each user
// @Tags urls
// @Produce json
// @Success 200 {object} models.GetURLsResponse
// @Failure 404 "Not Found"
// @Failue 500 {object} models.ErrorResponse
// @Router /shortener/urls [get]
func GetURLs(c *gin.Context) {
	// Retrieve all user objects from database
	userID, _ := c.Get("userID")
	var urls []models.URL
	config.DB.Where("user_id = ?", userID).Find(&urls)
	// Get count of users in database
	var count int64
	result := config.DB.Model(&models.URL{}).Where("user_id = ?", userID).Count(&count)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: "failed to get count"})
		return
	}
	// Return user details as JSON response
	c.JSON(http.StatusOK, models.GetURLsResponse{
		Success: true,
		Count:   count,
		URLs:    urls,
	})
}

// DeleteURL deletes an existing URL object from the database
//
// @Description Delete a url by its ID
// @Tags urls
// @Produce json
// @Param urlID path int true "Url ID"
// @Success 204 "No Content"
// @Failure 404 {object} models.ErrorResponse
// @Failue 500 {object} models.ErrorResponse
// @Router /shortener/url/{urlID} [delete]
func DeleteURL(c *gin.Context) {
	urlID := c.Param("urlID")
	userID, _ := c.Get("userID")
	var url models.URL

	config.DB.Where("id = ? AND user_id = ?", urlID, userID).First(&url)

	if !urlExists(urlID) {
		c.IndentedJSON(http.StatusNotFound, models.ErrorResponse{
			Error: fmt.Sprintf("URL object with id {%v} does not exist", urlID),
		})
		return
	}
	if err := config.DB.Unscoped().Delete(&url).Error; err != nil {
		// Handle deletion error
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: "failed to delete URL - you cannot delete a URL you did not create"})
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
