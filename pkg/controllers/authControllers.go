package controllers

import (
	"log"
	"net/http"
	"time"

	"github.com/0xMarvell/scissor/pkg/config"
	"github.com/0xMarvell/scissor/pkg/models"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"golang.org/x/crypto/bcrypt"
)

const JWT_SECRET = "secret"

// @Summary Create a new user
// @Description Create a new user with the provided details
// @Tags users
// @Accept json
// @Produce json
// @Success 201 {object} models.User
// @Failure 400 {object} string
// @Router /api/v1/signup [post]
func Signup(c *gin.Context) {
	// Get username, password off request body
	var signupPayload struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}
	if err := c.Bind(&signupPayload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to read request body",
		})
		return
	}

	// Hash password gotten from request body
	hash, err := bcrypt.GenerateFromPassword([]byte(signupPayload.Password), 10)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to hash password",
		})
		return
	}
	// Create user object using GORM
	user := models.User{Username: signupPayload.Username, Password: string(hash)}
	result := config.DB.Create(&user)
	if result.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to create user: an account with that username already exists",
		})
		log.Println(result.Error)
		return
	}
	// Return JSON response to confirm successful creation of user
	c.JSON(http.StatusCreated, gin.H{
		"success": true,
		"user":    user,
		"message": "New user was successfully created. Proceed to login",
	})
}

// Login allows existing user to login to the API
//
// @Summary Login to user account
// @Description Login with user details
// @Tags users
// @Accept json
// @Produce json
// @Success 200 {object} string
// @Failure 400 {object} string
// @Router /api/v1/login [post]
func Login(c *gin.Context) {
	// Get needed details (username,password) off request body
	var loginPayload struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	if c.Bind(&loginPayload) != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to read request body",
		})
		return
	}
	// Using GORM, query database to find user details
	var user models.User
	config.DB.First(&user, "username = ?", loginPayload.Username)
	if user.ID == 0 {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "User not found. Invalid username or password",
		})
		return
	}
	// Compare password gotten off request body to user password hash stored in database
	pwdErr := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(loginPayload.Password))
	if pwdErr != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid email or password",
		})
		return
	}
	// Generate JWT token
	// Create a new token object, specifying signing method and the claims
	// you would like it to contain.
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"subject": user.ID,
		/* Expiration time on line 94 == 30 days.
		   This is jut an example implementation
		   For production, 30 days would be too much so a shorter time
		   would be more optimal.
		*/
		"expiration_time": time.Now().Add(time.Hour * 24 * 30).Unix(),
	})
	// Sign and get the complete encoded token as a string using the secret
	// the secret key will be created by you (it can be a random sequence of characters e.g. 3r4jgnirbg8rhg08gvi0pvhh8)
	// Tip: DO NOT HARD CODE YOUR SECRET KEY in production
	// Store it as an environment variable instead
	tokenString, err := token.SignedString([]byte(JWT_SECRET))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to generate JWT token",
		})
		return
	}
	// Store JWT token inside httpOnly cookie (for security purposes)
	// Avoid storing your token in localstorage because it
	// becomes vulnerable to Cross-Site-Scripting (XSS) attack
	var secure, httpOnly bool = false, true // in production, secure should be set to true
	var maxAge int = 3600 * 24 * 30         // maxAge is the amount of time (IN SECONDS) the cookie will be valid for
	c.SetSameSite(http.SameSiteLaxMode)
	c.SetCookie("auth_token", tokenString, maxAge, "", "", secure, httpOnly)
	// Return JSON Response to confirm successful storage of JWT token
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "login successful",
	})
}

// Logout logs out the curret logged-in user
//
// @Summary Logout of user account
// @Description Logout of user account
// @Tags users
// @Accept json
// @Produce json
// @Success 201 {object} string
// @Failure 400 {object} string
// @Router /api/v1/logout [get]
func Logout(c *gin.Context) {
	// Remove the cookie containing the JWT authorization token
	//by setting its expiration time to a past value
	var secure, httpOnly bool = false, true
	c.SetCookie("auth_token", "", -1, "", "", secure, httpOnly)

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "logout successful",
	})
}

// GetUsers retrieves all existing users
//
// @Summary Get all users
// @Description Retrieve all existing user accounts
// @Tags users
// @Produce json
// @Success 200 {object} models.User
// @Failure 404 {object} string
// @Router /api/v1/users [get]
func GetUsers(c *gin.Context) {
	// Retrieve all user objects from database
	var users []models.User
	config.DB.Preload("URLs").Find(&users)

	// Get count of users in database
	var count int64
	result := config.DB.Model(&models.User{}).Count(&count)
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
		"users":   users,
	})
}

// GetUser retrieves an existing user's account details
//
// @Summary Get a user by ID
// @Description Retrieve an existing user
// @Tags users
// @Produce json
// @Success 200 {object} models.User
// @Failure 404 {object} string
// @Router /api/v1/users/profile [get]
func GetUser(c *gin.Context) {
	// Retrieve user details attached to request after passing through middleware
	user, _ := c.Get("user")
	// Return user details as JSON response
	c.JSON(http.StatusOK, gin.H{
		"success":      true,
		"user_details": user,
	})
}
