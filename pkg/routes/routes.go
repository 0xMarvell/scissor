package routes

import (
	_ "github.com/0xMarvell/scissor/docs"
	"github.com/0xMarvell/scissor/pkg/controllers"
	"github.com/0xMarvell/scissor/pkg/middleware"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// RegisterRoutes handles all routing for the API
func RegisterRoutes(r *gin.Engine) {
	// Swagger Docs
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// Authentication routes
	auth := r.Group("/api/v1")
	{
		// Index page...a simple greeting
		auth.GET("", controllers.SayHello) // won't require authentication
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
		shortener.POST("/shorten", middleware.RequireAuth, controllers.Shorten)
		shortener.GET("/urls", middleware.RequireAuth, controllers.GetURLs)
		shortener.GET("/redirect/:shortURL", middleware.RequireAuth, controllers.Redirect)
		shortener.DELETE("/url/:urlID", middleware.RequireAuth, controllers.DeleteURL)
	}
}
