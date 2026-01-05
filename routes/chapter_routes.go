package routes

import (
	"manara/controllers"
	middleware "manara/middlewares"

	"github.com/gin-gonic/gin"
)

func ChapterRoutes(router gin.IRouter) {
	chapters := router.Group("/chapters")
	chapters.Use(middleware.AuthMiddleware()) // All chapter routes require auth
	chapters.Use(middleware.RoleMiddleware("admin", "super_admin", "teacher"))
	{
		chapters.GET("", controllers.GetChapters)
		chapters.GET("/:id", controllers.GetChapter)
		chapters.POST("", controllers.CreateChapter)
		chapters.PUT("/:id", controllers.UpdateChapter)
		chapters.DELETE("/:id", controllers.DeleteChapter)
	}
}
