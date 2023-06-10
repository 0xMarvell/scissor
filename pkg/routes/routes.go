package routes

import (
	"github.com/0xMarvell/scissor/pkg/controllers"
	"github.com/0xMarvell/scissor/pkg/middleware"
	"github.com/gin-gonic/gin"
)

// RegisterRoutes handles all routing for the API
func RegisterRoutes(r *gin.Engine) {
	// Index page...a simple greeting
	r.GET("/", controllers.SayHello) // won't require authentication

	// Authentication routes
	auth := r.Group("/api/v1")
	{
		auth.POST("/signup", controllers.Signup)
		auth.POST("/login", controllers.Login)
		auth.GET("/logout", controllers.Logout)
	}

	// User routes
	users := r.Group("/api/v1/users")
	{
		users.GET("", controllers.GetUsers) // won't require authentication
		users.GET("/profile", middleware.RequireAuth, controllers.GetUser)
	}

	// URL shortener routes
	shortener := r.Group("api/v1/shortener")
	{
		shortener.GET("/urls", middleware.RequireAuth, controllers.GetURLs)
		shortener.DELETE("/url/:urlID", middleware.RequireAuth, controllers.DeleteURL)
		shortener.POST("/shorten", middleware.RequireAuth, controllers.Shorten)
		shortener.GET("/redirect/:shortURL", middleware.RequireAuth, controllers.Redirect)
	}

}
