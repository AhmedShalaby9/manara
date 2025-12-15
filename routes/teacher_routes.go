package routes

import (
	"manara/controllers"
	middleware "manara/middlewares"

	"github.com/gin-gonic/gin"
)

func TeacherRoutes(router *gin.Engine) {
	teachers := router.Group("/teachers")
	{
		teachers.GET("", controllers.GetTeachers)
		teachers.GET("/:id", controllers.GetTeacher)

		var adminOnly = teachers.Group("")
		adminOnly.Use(middleware.AuthMiddleware())
		adminOnly.Use(middleware.RoleMiddleware("admin", "super_admin"))
		{
			adminOnly.POST("", controllers.CreateTeacher)
			adminOnly.PUT("/:id", controllers.UpdateTeacher)
			adminOnly.DELETE("/:id", controllers.DeleteTeacher)
		}
	}
}
