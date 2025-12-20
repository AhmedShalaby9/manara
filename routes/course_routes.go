package routes

import (
	"manara/controllers"
	middleware "manara/middlewares"

	"github.com/gin-gonic/gin"
)

func CourseRoutes(router gin.IRouter) { // ← Changed
	courses := router.Group("/courses")
	{
		courses.GET("", controllers.GetCourses)
		courses.GET("/:id", controllers.GetCourse)
		courses.GET("/teacher/:teacher_id", controllers.GetTeacherCourses)

		adminOnly := courses.Group("")
		adminOnly.Use(middleware.AuthMiddleware())
		adminOnly.Use(middleware.RoleMiddleware("admin", "super_admin"))
		{
			adminOnly.POST("", controllers.CreateCourse)
			adminOnly.PUT("/:id", controllers.UpdateCourse)
			adminOnly.DELETE("/:id", controllers.DeleteCourse)
			adminOnly.POST("/assign", controllers.AssignCourseToTeacher)
			adminOnly.POST("/:id/image", controllers.UploadCourseImage)
			adminOnly.DELETE("/:id/image", controllers.DeleteCourseImage)
		}
	}
}
