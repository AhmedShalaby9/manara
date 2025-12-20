package routes

import (
	"manara/controllers"

	"github.com/gin-gonic/gin"
)

func UserRoutes(router gin.IRouter) {
	users := router.Group("/users")
	{
		users.GET("", controllers.GetUsers)
		users.POST("", controllers.CreateUser)
		users.PUT("/:id", controllers.ToggleUserActivation)

	}
}
