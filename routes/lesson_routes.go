package routes

import (
	"manara/controllers"
	middleware "manara/middlewares"

	"github.com/gin-gonic/gin"
)

func LessonRoutes(router gin.IRouter) {
	lessons := router.Group("/lessons")
	lessons.Use(middleware.AuthMiddleware())
	{
		// Read access for all authenticated users (admin, teacher, student)
		readAccess := lessons.Group("")
		readAccess.Use(middleware.RoleMiddleware("admin", "super_admin", "teacher", "student"))
		{
			readAccess.GET("", controllers.GetLessons)
			readAccess.GET("/:id", controllers.GetLesson)
			readAccess.GET("/:id/files", controllers.GetLessonFiles)
			readAccess.GET("/:id/videos", controllers.GetLessonVideos)
		}

		// Write access for admin and teacher only
		writeAccess := lessons.Group("")
		writeAccess.Use(middleware.RoleMiddleware("admin", "super_admin", "teacher"))
		{
			writeAccess.POST("", controllers.CreateLesson)
			writeAccess.PUT("/:id", controllers.UpdateLesson)
			writeAccess.DELETE("/:id", controllers.DeleteLesson)

			// File upload endpoints
			writeAccess.POST("/:id/files", controllers.UploadLessonFiles)
			writeAccess.DELETE("/:id/files/:file_id", controllers.DeleteLessonFile)

			// Video upload endpoints
			writeAccess.POST("/:id/videos", controllers.UploadLessonVideos)
			writeAccess.DELETE("/:id/videos/:video_id", controllers.DeleteLessonVideo)
		}
	}
}
