package routes

import (
	"manara/controllers"
	middleware "manara/middlewares"

	"github.com/gin-gonic/gin"
)

func ChapterRoutes(router gin.IRouter) {
	chapters := router.Group("/chapters")
	chapters.Use(middleware.AuthMiddleware())
	{
		// Read access for all authenticated users (admin, teacher, student)
		readAccess := chapters.Group("")
		readAccess.Use(middleware.RoleMiddleware("admin", "super_admin", "teacher", "student"))
		{
			readAccess.GET("", controllers.GetChapters)
			readAccess.GET("/:id", controllers.GetChapter)
		}

		// Write access for admin and teacher only
		writeAccess := chapters.Group("")
		writeAccess.Use(middleware.RoleMiddleware("admin", "super_admin", "teacher"))
		{
			writeAccess.POST("", controllers.CreateChapter)
			writeAccess.PUT("/:id", controllers.UpdateChapter)
			writeAccess.DELETE("/:id", controllers.DeleteChapter)
		}
	}
}
