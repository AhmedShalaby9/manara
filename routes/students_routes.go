package routes

import (
	"manara/controllers"

	"github.com/gin-gonic/gin"
)

func StudentRoutes(router *gin.Engine) {
	students := router.Group("/students")
	{
		students.GET("", controllers.GetStudents)
		students.GET("/:id", controllers.GetStudent)
		students.POST("", controllers.CreateStudent)

	}
}
