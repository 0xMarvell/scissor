package routes

import (
	"github.com/0xMarvell/scissor/pkg/controllers"
	"github.com/0xMarvell/scissor/pkg/middleware"
	"github.com/gin-gonic/gin"
)

// RegisterUserRoutes handles all routing for the API
func RegisterRoutes(r *gin.Engine) {
	// Authentication routes
	auth := r.Group("/api/v1/user")
	{
		auth.POST("/signup", controllers.Signup)
		auth.POST("/login", controllers.Login)
		auth.GET("/logout", controllers.Logout)
		auth.GET("/details", middleware.RequireAuth, controllers.GetUserDetails)
	}

	r.GET("/", middleware.RequireAuth, controllers.SayHello)

	// URL shortener routes
	// shortener := r.Group("api/v1/shortener")
	// {

	// }

}
