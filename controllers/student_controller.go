package controllers

import (
	"strconv"

	"manara/database"
	"manara/helpers"
	"manara/models"

	"github.com/gin-gonic/gin"
)

func GetStudents(c *gin.Context) {
	var students []models.Student
	teacherID := c.Query("teacher_id")
	gradeLevel := c.Query("grade_level")

	query := database.DB

	if teacherID != "" {
		query = query.Where("teacher_id = ?", teacherID)
	}
	if gradeLevel != "" {
		query = query.Where("grade_level =?", gradeLevel)
	}

	res := query.Find(&students)
	if res.Error != nil {
		helpers.Respond(c, false, nil, res.Error.Error())
		return
	}

	helpers.Respond(c, true, students, "Students retrieved successfully")
}

func GetStudent(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	var student models.Student

	res := database.DB.Preload("User.Role").First(&student, id)
	if res.Error != nil {
		helpers.Respond(c, false, nil, "Student not found")
		return
	}

	helpers.Respond(c, true, student, "Student retrieved successfully")
}

func CreateStudent(c *gin.Context) {
	var req models.CreateStudentRequest

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
		RoleID:       4,
		IsActive:     true,
	}

	if err := tx.Create(&user).Error; err != nil {
		tx.Rollback()
		helpers.Respond(c, false, nil, "Failed to create user")
		return
	}

	student := models.Student{
		UserID:    user.ID,
		TeacherID: req.TeacherId, GradeLevel: req.GradeLevel, ParentPhone: req.ParentPhone,
	}

	if err := tx.Create(&student).Error; err != nil {
		tx.Rollback()
		helpers.Respond(c, false, nil, "Failed to create Student")
		return
	}

	tx.Commit()

	database.DB.Preload("User.Role").First(&student, student.ID)

	helpers.Respond(c, true, student, "Student created successfully")
}

// UpdateTeacher - Update teacher info
// func UpdateStudent(c *gin.Context) {
// 	id, _ := strconv.Atoi(c.Param("id"))
// 	var teacher models.Teacher

// 	// Check if teacher exists
// 	if err := database.DB.First(&teacher, id).Error; err != nil {
// 		helpers.Respond(c, false, nil, "Teacher not found")
// 		return
// 	}

// 	var req models.UpdateTeacherRequest
// 	if err := c.ShouldBindJSON(&req); err != nil {
// 		helpers.Respond(c, false, nil, err.Error())
// 		return
// 	}

// 	// Update teacher info
// 	teacher.Bio = req.Bio
// 	teacher.Specialization = req.Specialization
// 	teacher.ExperienceYears = req.ExperienceYears

// 	if err := database.DB.Save(&teacher).Error; err != nil {
// 		helpers.Respond(c, false, nil, "Failed to update teacher")
// 		return
// 	}

// 	// Load with user info
// 	database.DB.Preload("User.Role").First(&teacher, teacher.ID)

// 	helpers.Respond(c, true, teacher, "Teacher updated successfully")
// }

// DeleteTeacher - Delete a teacher (Admin/Super Admin only)
// func DeleteStudent(c *gin.Context) {
// 	id, _ := strconv.Atoi(c.Param("id"))

// 	var teacher models.Teacher
// 	if err := database.DB.First(&teacher, id).Error; err != nil {
// 		helpers.Respond(c, false, nil, "Teacher not found")
// 		return
// 	}

// 	// Start transaction
// 	tx := database.DB.Begin()

// 	// Delete teacher record
// 	if err := tx.Delete(&teacher).Error; err != nil {
// 		tx.Rollback()
// 		helpers.Respond(c, false, nil, "Failed to delete teacher")
// 		return
// 	}

// 	// Delete associated user
// 	if err := tx.Delete(&models.User{}, teacher.UserID).Error; err != nil {
// 		tx.Rollback()
// 		helpers.Respond(c, false, nil, "Failed to delete user")
// 		return
// 	}

// 	tx.Commit()

// 	helpers.Respond(c, true, nil, "Teacher deleted successfully")
// }
