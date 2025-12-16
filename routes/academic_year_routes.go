package routes

import (
	"manara/controllers"

	"github.com/gin-gonic/gin"
)

func AcademicYearRoutes(router *gin.Engine) {
	years := router.Group("/academicYears")
	{
		years.GET("", controllers.GetAcademicYears)
		years.GET("/:id", controllers.GetAcademicYear)
		years.POST("", controllers.CreateAcademicYear)
		years.DELETE("/:id", controllers.DeleteAcademicYear)

	}
}
