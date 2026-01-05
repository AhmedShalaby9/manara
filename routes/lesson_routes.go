package routes

import (
	"manara/controllers"
	middleware "manara/middlewares"

	"github.com/gin-gonic/gin"
)

func LessonRoutes(router gin.IRouter) {
	lessons := router.Group("/lessons")
	lessons.Use(middleware.AuthMiddleware()) // All lesson routes require auth
	lessons.Use(middleware.RoleMiddleware("admin", "super_admin", "teacher"))
	{
		lessons.GET("", controllers.GetLessons)
		lessons.GET("/:id", controllers.GetLesson)
		lessons.GET("/:id/files", controllers.GetLessonFiles)
		lessons.GET("/:id/videos", controllers.GetLessonVideos)
		lessons.POST("", controllers.CreateLesson)
		lessons.PUT("/:id", controllers.UpdateLesson)
		lessons.DELETE("/:id", controllers.DeleteLesson)

		// File upload endpoints
		lessons.POST("/:id/files", controllers.UploadLessonFiles)
		lessons.DELETE("/:id/files/:file_id", controllers.DeleteLessonFile)

		// Video upload endpoints
		lessons.POST("/:id/videos", controllers.UploadLessonVideos)
		lessons.DELETE("/:id/videos/:video_id", controllers.DeleteLessonVideo)
	}
}
