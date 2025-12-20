package routes

import (
	"manara/controllers"
	middleware "manara/middlewares"

	"github.com/gin-gonic/gin"
)

func AuthRoutes(router gin.IRouter) {
	auth := router.Group("/auth")
	{
		auth.POST("/register", controllers.Register)
		auth.POST("/login", controllers.Login)
		auth.PATCH("/change-password/:id", controllers.ChangePassword)
		auth.GET("/me", middleware.AuthMiddleware(), controllers.GetMe)
		auth.POST("/logout", middleware.AuthMiddleware(), controllers.Logout)
		auth.GET("/generate-username", controllers.GenerateUniqueUserName)
	}
}
