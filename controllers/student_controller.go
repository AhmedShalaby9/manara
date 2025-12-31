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
	academicYearID := c.Query("academic_year_id")
	search := c.Query("search")

	params := helpers.GetPaginationParams(c)
	query := database.DB.Model(&models.Student{}).Preload("User").Preload("Teacher.User").Preload("AcademicYear")

	if teacherID != "" {
		query = query.Where("teacher_id = ?", teacherID)
	}
	if gradeLevel != "" {
		query = query.Where("grade_level = ?", gradeLevel)
	}
	if academicYearID != "" {
		query = query.Where("academic_year_id = ?", academicYearID)
	}
	if search != "" {
		query = query.Joins("JOIN users ON users.id = students.user_id").
			Where("users.first_name LIKE ? OR users.last_name LIKE ? OR users.email LIKE ?", "%"+search+"%", "%"+search+"%", "%"+search+"%")
	}

	pagination, err := helpers.Paginate(query, params, &students)
	if err != nil {
		helpers.Respond(c, false, nil, "Failed to retrieve students")
		return
	}

	helpers.RespondWithPagin(c, true, students, "Students retrieved successfully", pagination)
}

func GetStudent(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	var student models.Student

	res := database.DB.Preload("User.Role").Preload("Teacher.User").Preload("AcademicYear").First(&student, id)
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

	var academicYear models.AcademicYear
	if err := database.DB.First(&academicYear, req.AcademicYearID).Error; err != nil {
		helpers.Respond(c, false, nil, "Academic year is wrong!")
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
		UserID:         user.ID,
		TeacherID:      req.TeacherID,
		AcademicYearID: req.AcademicYearID,
		ParentPhone:    req.ParentPhone,
	}

	if err := tx.Create(&student).Error; err != nil {
		tx.Rollback()
		helpers.Respond(c, false, nil, "Failed to create Student")
		return
	}

	tx.Commit()

	database.DB.Preload("User").Preload("AcademicYear").First(&student, student.ID)

	helpers.Respond(c, true, student, "Student created successfully")
}
