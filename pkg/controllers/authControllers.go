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

// @Description Create a new user with the provided details
// @Tags authentication
// @Accept json
// @Produce json
// @Param payload body models.SignupRequest true "payload"
// @Success 201 {object} models.SignupResponse
// @Failure 400 {object} models.ErrorResponse
// @Router /signup [post]
func Signup(c *gin.Context) {
	// Get username, password off request body
	var s models.SignupRequest
	if err := c.Bind(&s); err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error: "Failed to read request body",
		})
		return
	}

	// Hash password gotten from request body
	hash, err := bcrypt.GenerateFromPassword([]byte(s.Password), 10)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error: "Failed to hash password",
		})
		return
	}
	// Create user object using GORM
	user := models.User{Username: s.Username, Password: string(hash)}
	result := config.DB.Create(&user)
	if result.Error != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error: "Failed to create user: an account with that username already exists",
		})
		log.Println(result.Error)
		return
	}
	// Return JSON response to confirm successful creation of user
	c.JSON(http.StatusCreated, models.SignupResponse{
		Success: true,
		Message: "New user was successfully created. Proceed to login",
		User:    user,
	})
}

// Login allows existing user to login to the API
//
// @Description Login with user details
// @Tags authentication
// @Accept json
// @Produce json
// @Param payload body models.LoginRequest true "payload"
// @Success 200 {object} models.LoginResponse
// @Failure 400 {object} models.ErrorResponse
// @Failure 404 {object} models.ErrorResponse
// @Router /login [post]
func Login(c *gin.Context) {
	// Get needed details (username,password) off request body
	var l models.LoginRequest

	if c.Bind(&l) != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error: "Failed to read request body",
		})
		return
	}
	// Using GORM, query database to find user details
	var user models.User
	config.DB.First(&user, "username = ?", l.Username)
	if user.ID == 0 {
		c.JSON(http.StatusNotFound, models.ErrorResponse{
			Error: "User not found. Invalid username or password",
		})
		return
	}
	// Compare password gotten off request body to user password hash stored in database
	pwdErr := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(l.Password))
	if pwdErr != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error: "Invalid email or password",
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
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error: "Failed to generate JWT token",
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
	c.JSON(http.StatusOK, models.LoginResponse{
		Success: true,
		Message: "login successful",
	})
}

// Logout logs out the curret logged-in user
//
// @Description Logout of user account
// @Tags authentication
// @Produce json
// @Success 200 {object} models.LogoutResponse
// @Failure 400 "Bad Request"
// @Router /logout [get]
func Logout(c *gin.Context) {
	// Remove the cookie containing the JWT authorization token
	//by setting its expiration time to a past value
	var secure, httpOnly bool = false, true
	c.SetCookie("auth_token", "", -1, "", "", secure, httpOnly)

	c.JSON(http.StatusOK, models.LogoutResponse{
		Success: true,
		Message: "logout successful",
	})
}

// GetUsers retrieves all existing users
//
// @Description Retrieve all existing user accounts
// @Tags users
// @Produce json
// @Success 200 {object} models.GetUsersResponse
// @Failure 404 "Not Found"
// @Router /users [get]
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
	c.JSON(http.StatusOK, models.GetUsersResponse{
		Success: true,
		Count:   count,
		Users:   users,
	})
}

// GetUser retrieves an existing user's account details
//
// @Description Retrieve an existing user
// @Tags users
// @Produce json
// @Success 200 {object} models.GetUserResponse
// @Failure 404 "Not Found"
// @Failure 401 {object} models.ErrorResponse
// @Router /users/profile [get]
func GetUser(c *gin.Context) {
	// Retrieve user details attached to request after passing through middleware
	user, _ := c.Get("user")
	// Return user details as JSON response
	c.JSON(http.StatusOK, models.GetUserResponse{
		Success: true,
		User:    user,
	})
}
