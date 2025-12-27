package routes

import (
	"manara/controllers"
	middleware "manara/middlewares"

	"github.com/gin-gonic/gin"
)

func UserRoutes(router gin.IRouter) {
	users := router.Group("/users")
	{

		adminOnly := users.Group("")
		adminOnly.Use(middleware.AuthMiddleware())
		adminOnly.Use(middleware.RoleMiddleware("admin", "super_admin"))

		adminOnly.GET("", controllers.GetUsers)
		adminOnly.POST("", controllers.CreateUser)
		adminOnly.PUT("/:id", controllers.ToggleUserActivation)
		adminOnly.DELETE("/:id", controllers.DeleteUser)

	}
}
