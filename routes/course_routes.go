package routes

import (
	"manara/controllers"
	middleware "manara/middlewares"

	"github.com/gin-gonic/gin"
)

func CourseRoutes(router gin.IRouter) {
	courses := router.Group("/courses")
	{
		courses.GET("", controllers.GetCourses)
		courses.GET("/:id", controllers.GetCourse)

		// Get my courses (teacher sees their courses, student sees teacher's courses)
		authRoutes := courses.Group("")
		authRoutes.Use(middleware.AuthMiddleware())
		authRoutes.Use(middleware.RoleMiddleware("teacher", "student", "admin", "super_admin"))
		{
			authRoutes.GET("/my", controllers.GetMyCourses)
		}

		// Admin only routes
		adminOnly := courses.Group("")
		adminOnly.Use(middleware.AuthMiddleware())
		adminOnly.Use(middleware.RoleMiddleware("admin", "super_admin"))
		{
			adminOnly.GET("/teacher/:teacher_id", controllers.GetTeacherCourses)
			adminOnly.POST("", controllers.CreateCourse)
			adminOnly.PUT("/:id", controllers.UpdateCourse)
			adminOnly.DELETE("/:id", controllers.DeleteCourse)
			adminOnly.POST("/assign", controllers.AssignCourseToTeacher)
			adminOnly.POST("/:id/image", controllers.UploadCourseImage)
			adminOnly.DELETE("/:id/image", controllers.DeleteCourseImage)
		}
	}
}
