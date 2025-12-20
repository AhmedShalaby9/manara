package routes

import (
	"manara/controllers"
	middleware "manara/middlewares"

	"github.com/gin-gonic/gin"
)

func ChapterRoutes(router gin.IRouter) { // ← Changed
	chapters := router.Group("/chapters")
	{
		chapters.GET("", controllers.GetChapters)
		chapters.GET("/:id", controllers.GetChapter)

		teacherAndAdmin := chapters.Group("")
		teacherAndAdmin.Use(middleware.AuthMiddleware())
		teacherAndAdmin.Use(middleware.RoleMiddleware("admin", "super_admin", "teacher"))
		{
			teacherAndAdmin.POST("", controllers.CreateChapter)
			teacherAndAdmin.PUT("/:id", controllers.UpdateChapter)
			teacherAndAdmin.DELETE("/:id", controllers.DeleteChapter)
		}
	}
}
