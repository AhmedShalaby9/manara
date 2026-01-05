package routes

import (
	"manara/controllers"
	middleware "manara/middlewares"

	"github.com/gin-gonic/gin"
)

func StudentRoutes(router gin.IRouter) {
	students := router.Group("/students")
	students.Use(middleware.AuthMiddleware()) // All student routes require auth
	students.Use(middleware.RoleMiddleware("admin", "super_admin", "teacher"))
	{
		students.GET("", controllers.GetStudents)
		students.GET("/:id", controllers.GetStudent)
		students.POST("", controllers.CreateStudent)
	}
}
