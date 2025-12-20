package routes

import (
	"manara/controllers"
	middleware "manara/middlewares"

	"github.com/gin-gonic/gin"
)

func LessonRoutes(router gin.IRouter) { // ← Changed
	lessons := router.Group("/lessons")
	{
		lessons.GET("", controllers.GetLessons)
		lessons.GET("/:id", controllers.GetLesson)

		teacherAndAdmin := lessons.Group("")
		teacherAndAdmin.Use(middleware.AuthMiddleware())
		teacherAndAdmin.Use(middleware.RoleMiddleware("admin", "super_admin", "teacher"))
		{
			teacherAndAdmin.POST("", controllers.CreateLesson)
			teacherAndAdmin.PUT("/:id", controllers.UpdateLesson)
			teacherAndAdmin.DELETE("/:id", controllers.DeleteLesson)
		}
	}
}
