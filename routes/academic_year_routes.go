package routes

import (
	"manara/controllers"
	middleware "manara/middlewares"

	"github.com/gin-gonic/gin"
)

func AcademicYearRoutes(router gin.IRouter) { // ← Changed
	academicYears := router.Group("/academicYears")
	{
		academicYears.GET("", controllers.GetAcademicYears)
		academicYears.GET("/:id", controllers.GetAcademicYear)

		adminOnly := academicYears.Group("")
		adminOnly.Use(middleware.AuthMiddleware())
		adminOnly.Use(middleware.RoleMiddleware("admin", "super_admin"))
		{
			adminOnly.POST("", controllers.CreateAcademicYear)
			adminOnly.DELETE("/:id", controllers.DeleteAcademicYear)
		}
	}
}
