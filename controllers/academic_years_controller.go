package controllers

import (
	"manara/database"
	"manara/helpers"
	"manara/models"
	"strconv"

	"github.com/gin-gonic/gin"
)

func GetAcademicYears(c *gin.Context) {
	var academicYear []models.AcademicYear
	res := database.DB.Find(&academicYear)
	if res.Error != nil {
		helpers.Respond(c, false, nil, res.Error.Error())
		return
	}
	helpers.Respond(c, true, academicYear, "Academic Years retrieved successfully")
}

func GetAcademicYear(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	var academicYear models.AcademicYear
	res := database.DB.First(&academicYear, id)
	if res.Error != nil {
		helpers.Respond(c, false, nil, "Academic Year not found")
		return
	}
	helpers.Respond(c, true, academicYear, "Academic Year retrieved successfully")
}
func CreateAcademicYear(c *gin.Context) {
	var academicYear models.AcademicYear
	if err := c.ShouldBindJSON(&academicYear); err != nil {
		helpers.Respond(c, false, nil, err.Error())
		return
	}
	if academicYear.Name == "" {
		helpers.Respond(c, false, nil, "Academic year name is required")
		return
	}

	res := database.DB.Create(&academicYear)
	if res.Error != nil {
		helpers.Respond(c, false, nil, res.Error.Error())
		return
	}
	helpers.Respond(c, true, academicYear, "Role created successfully")
}

func DeleteAcademicYear(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))

	var year models.AcademicYear
	if err := database.DB.First(&year, id).Error; err != nil {
		helpers.Respond(c, false, nil, "Academic year not found")
		return
	}

	res := database.DB.Delete(&year, id)
	if res.Error != nil {
		helpers.Respond(c, false, nil, res.Error.Error())
		return
	}
	helpers.Respond(c, true, nil, "Academic year deleted successfully")
}
