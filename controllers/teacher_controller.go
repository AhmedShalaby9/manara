package controllers

import (
	"strconv"

	"manara/database"
	"manara/helpers"
	"manara/models"

	"github.com/gin-gonic/gin"
)

func GetTeachers(c *gin.Context) {
	var teachers []models.Teacher

	res := database.DB.Preload("User.Role").Find(&teachers)
	if res.Error != nil {
		helpers.Respond(c, false, nil, res.Error.Error())
		return
	}

	helpers.Respond(c, true, teachers, "Teachers retrieved successfully")
}

func GetTeacher(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	var teacher models.Teacher

	res := database.DB.Preload("User.Role").First(&teacher, id)
	if res.Error != nil {
		helpers.Respond(c, false, nil, "Teacher not found")
		return
	}

	helpers.Respond(c, true, teacher, "Teacher retrieved successfully")
}

func CreateTeacher(c *gin.Context) {
	var req models.CreateTeacherRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		helpers.Respond(c, false, nil, err.Error())
		return
	}

	var existingUser models.User
	if err := database.DB.Where("user_name = ?", req.UserName).First(&existingUser).Error; err == nil {
		helpers.Respond(c, false, nil, "Username already exists")
		return
	}

	if err := database.DB.Where("email = ?", req.Email).First(&existingUser).Error; err == nil {
		helpers.Respond(c, false, nil, "Email already exists")
		return
	}

	hashedPassword, err := helpers.HashPassword(req.Password)
	if err != nil {
		helpers.Respond(c, false, nil, "Failed to hash password")
		return
	}

	tx := database.DB.Begin()

	user := models.User{
		FirstName:    req.FirstName,
		LastName:     req.LastName,
		Email:        req.Email,
		Phone:        req.Phone,
		UserName:     req.UserName,
		PasswordHash: hashedPassword,
		RoleID:       3,
		IsActive:     true,
	}

	if err := tx.Create(&user).Error; err != nil {
		tx.Rollback()
		helpers.Respond(c, false, nil, "Failed to create user")
		return
	}

	teacher := models.Teacher{
		UserID:          user.ID,
		Bio:             req.Bio,
		Specialization:  req.Specialization,
		ExperienceYears: req.ExperienceYears,
	}

	if err := tx.Create(&teacher).Error; err != nil {
		tx.Rollback()
		helpers.Respond(c, false, nil, "Failed to create teacher")
		return
	}

	tx.Commit()

	database.DB.Preload("User.Role").First(&teacher, teacher.ID)

	helpers.Respond(c, true, teacher, "Teacher created successfully")
}

func UpdateTeacher(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	var teacher models.Teacher

	if err := database.DB.First(&teacher, id).Error; err != nil {
		helpers.Respond(c, false, nil, "Teacher not found")
		return
	}

	var req models.UpdateTeacherRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		helpers.Respond(c, false, nil, err.Error())
		return
	}

	teacher.Bio = req.Bio
	teacher.Specialization = req.Specialization
	teacher.ExperienceYears = req.ExperienceYears

	if err := database.DB.Save(&teacher).Error; err != nil {
		helpers.Respond(c, false, nil, "Failed to update teacher")
		return
	}

	// Load with user info
	database.DB.Preload("User.Role").First(&teacher, teacher.ID)

	helpers.Respond(c, true, teacher, "Teacher updated successfully")
}

// DeleteTeacher - Delete a teacher (Admin/Super Admin only)
func DeleteTeacher(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))

	var teacher models.Teacher
	if err := database.DB.First(&teacher, id).Error; err != nil {
		helpers.Respond(c, false, nil, "Teacher not found")
		return
	}

	// Start transaction
	tx := database.DB.Begin()

	// Delete teacher record
	if err := tx.Delete(&teacher).Error; err != nil {
		tx.Rollback()
		helpers.Respond(c, false, nil, "Failed to delete teacher")
		return
	}

	// Delete associated user
	if err := tx.Delete(&models.User{}, teacher.UserID).Error; err != nil {
		tx.Rollback()
		helpers.Respond(c, false, nil, "Failed to delete user")
		return
	}

	tx.Commit()

	helpers.Respond(c, true, nil, "Teacher deleted successfully")
}
