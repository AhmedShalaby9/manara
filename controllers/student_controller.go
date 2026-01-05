package controllers

import (
	"strconv"

	"manara/database"
	"manara/helpers"
	"manara/models"

	"github.com/gin-gonic/gin"
)

// GetStudents - Get all students
// For teachers: automatically filtered to their own students
// For admins: can filter by teacher_id query param or see all
func GetStudents(c *gin.Context) {
	var students []models.Student
	gradeLevel := c.Query("grade_level")
	academicYearID := c.Query("academic_year_id")
	search := c.Query("search")

	params := helpers.GetPaginationParams(c)
	query := database.DB.Model(&models.Student{}).Preload("User").Preload("Teacher.User").Preload("AcademicYear")

	// Role-based teacher scoping
	if teacherID := helpers.GetEffectiveTeacherID(c); teacherID != nil {
		query = query.Where("teacher_id = ?", *teacherID)
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

// CreateStudent - Create a new student
// For teachers: teacher_id is automatically set from token
// For admins: teacher_id must be provided in request
func CreateStudent(c *gin.Context) {
	var req models.CreateStudentRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		helpers.Respond(c, false, nil, err.Error())
		return
	}

	// Get effective teacher_id based on role
	teacherID, ok := helpers.GetTeacherIDForCreate(c, req.TeacherID)
	if !ok {
		helpers.Respond(c, false, nil, "Teacher ID is required")
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
		TeacherID:      teacherID,
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
