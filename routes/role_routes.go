package routes

import (
	"manara/controllers"
	middleware "manara/middlewares"

	"github.com/gin-gonic/gin"
)

func RoleRoutes(router gin.IRouter) { // ← Changed
	roles := router.Group("/roles")
	{
		roles.GET("", controllers.GetRoles)
		roles.GET("/:id", controllers.GetRole)

		adminOnly := roles.Group("")
		adminOnly.Use(middleware.AuthMiddleware())
		adminOnly.Use(middleware.RoleMiddleware("admin", "super_admin"))
		{
			adminOnly.POST("", controllers.CreateRole)
			adminOnly.PUT("/:id", controllers.UpdateRole)
			adminOnly.DELETE("/:id", controllers.DeleteRole)
		}
	}
}
