package routes

import (
	"manara/controllers"
	middleware "manara/middlewares"

	"github.com/gin-gonic/gin"
)

func LessonRoutes(router gin.IRouter) {
	lessons := router.Group("/lessons")
	{
		lessons.GET("", controllers.GetLessons)
		lessons.GET("/:id", controllers.GetLesson)
		lessons.GET("/:id/files", controllers.GetLessonFiles)
		lessons.GET("/:id/videos", controllers.GetLessonVideos)

		teacherAndAdmin := lessons.Group("")
		teacherAndAdmin.Use(middleware.AuthMiddleware())
		teacherAndAdmin.Use(middleware.RoleMiddleware("admin", "super_admin", "teacher"))
		{
			teacherAndAdmin.POST("", controllers.CreateLesson)
			teacherAndAdmin.PUT("/:id", controllers.UpdateLesson)
			teacherAndAdmin.DELETE("/:id", controllers.DeleteLesson)

			// File upload endpoints
			teacherAndAdmin.POST("/:id/files", controllers.UploadLessonFiles)
			teacherAndAdmin.DELETE("/:id/files/:file_id", controllers.DeleteLessonFile)

			// Video upload endpoints
			teacherAndAdmin.POST("/:id/videos", controllers.UploadLessonVideos)
			teacherAndAdmin.DELETE("/:id/videos/:video_id", controllers.DeleteLessonVideo)
		}
	}
}
