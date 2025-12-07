package controllers

import (
	"manara/database"
	"manara/helpers"
	"manara/models"
	"strconv"

	"github.com/gin-gonic/gin"
)

func GetRoles(c *gin.Context) {
	var roles []models.Role
	res := database.DB.Find(&roles)
	if res.Error != nil {
		helpers.Respond(c, false, nil, res.Error.Error())
		return
	}
	helpers.Respond(c, true, roles, "Roles retrieved successfully")
}

func GetRole(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	var role models.Role
	res := database.DB.First(&role, id)
	if res.Error != nil {
		helpers.Respond(c, false, nil, "Role not found")
		return
	}
	helpers.Respond(c, true, role, "Role retrieved successfully")
}

func CreateRole(c *gin.Context) {
	var role models.Role
	if err := c.ShouldBindJSON(&role); err != nil {
		helpers.Respond(c, false, nil, err.Error())
		return
	}

	if role.RoleName == "" {
		helpers.Respond(c, false, nil, "Role name is required")
		return
	}
	if role.RoleValue == "" {
		helpers.Respond(c, false, nil, "Role value is required")
		return
	}

	res := database.DB.Create(&role)
	if res.Error != nil {
		helpers.Respond(c, false, nil, res.Error.Error())
		return
	}
	helpers.Respond(c, true, role, "Role created successfully")
}

func UpdateRole(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	var role models.Role

	if err := database.DB.First(&role, id).Error; err != nil {
		helpers.Respond(c, false, nil, "Role not found")
		return
	}

	if err := c.ShouldBindJSON(&role); err != nil {
		helpers.Respond(c, false, nil, err.Error())
		return
	}

	if role.RoleName == "" {
		helpers.Respond(c, false, nil, "Role name is required")
		return
	}
	if role.RoleValue == "" {
		helpers.Respond(c, false, nil, "Role value is required")
		return
	}

	res := database.DB.Save(&role)
	if res.Error != nil {
		helpers.Respond(c, false, nil, res.Error.Error())
		return
	}
	helpers.Respond(c, true, role, "Role updated successfully")
}

func DeleteRole(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))

	// Check if role exists
	var role models.Role
	if err := database.DB.First(&role, id).Error; err != nil {
		helpers.Respond(c, false, nil, "Role not found")
		return
	}

	res := database.DB.Delete(&role, id)
	if res.Error != nil {
		helpers.Respond(c, false, nil, res.Error.Error())
		return
	}
	helpers.Respond(c, true, nil, "Role deleted successfully")
}
