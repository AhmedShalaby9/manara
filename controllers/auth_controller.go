package controllers

import (
	"manara/database"
	"manara/helpers"
	"manara/models"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

// RegisterStudent - Register a new student (for mobile app)
// Student provides their teacher's username to link to the correct teacher
func RegisterStudent(c *gin.Context) {
	var req models.RegisterStudentRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		helpers.Respond(c, false, nil, err.Error())
		return
	}

	// Find teacher by username
	var teacherUser models.User
	if err := database.DB.Where("user_name = ?", req.TeacherUserName).First(&teacherUser).Error; err != nil {
		helpers.Respond(c, false, nil, "Teacher not found")
		return
	}

	// Get the teacher record
	var teacher models.Teacher
	if err := database.DB.Where("user_id = ?", teacherUser.ID).First(&teacher).Error; err != nil {
		helpers.Respond(c, false, nil, "Teacher not found")
		return
	}

	// Check if username already exists
	var existingUser models.User
	if err := database.DB.Where("user_name = ?", req.UserName).First(&existingUser).Error; err == nil {
		helpers.Respond(c, false, nil, "Username already exists")
		return
	}

	// Check if email already exists
	if err := database.DB.Where("email = ?", req.Email).First(&existingUser).Error; err == nil {
		helpers.Respond(c, false, nil, "Email already exists")
		return
	}

	// Hash password
	hashedPassword, err := helpers.HashPassword(req.Password)
	if err != nil {
		helpers.Respond(c, false, nil, "Failed to hash password")
		return
	}

	// Get student role by value
	var studentRole models.Role
	if err := database.DB.Where("role_value = ?", "student").First(&studentRole).Error; err != nil {
		helpers.Respond(c, false, nil, "Student role not found")
		return
	}

	// Start transaction
	tx := database.DB.Begin()

	// Create user with student role
	user := models.User{
		FirstName:    req.FirstName,
		LastName:     req.LastName,
		Email:        req.Email,
		Phone:        req.Phone,
		UserName:     req.UserName,
		PasswordHash: hashedPassword,
		RoleID:       studentRole.ID,
		IsActive:     false,
	}

	if err := tx.Create(&user).Error; err != nil {
		tx.Rollback()
		helpers.Respond(c, false, nil, "Failed to create user")
		return
	}

	// Create student record linked to teacher
	student := models.Student{
		UserID:      user.ID,
		TeacherID:   teacher.ID,
		ParentPhone: req.ParentPhone,
	}

	if err := tx.Create(&student).Error; err != nil {
		tx.Rollback()
		helpers.Respond(c, false, nil, "Failed to create student")
		return
	}

	tx.Commit()

	// Load user with role
	database.DB.Preload("Role").First(&user, user.ID)

	// Generate token
	token, err := helpers.GenerateToken(user.ID, user.UserName, user.RoleID, "student", nil)
	if err != nil {
		helpers.Respond(c, false, nil, "Failed to generate token")
		return
	}

	response := models.AuthResponse{
		Token: token,
		User:  user,
	}

	c.JSON(http.StatusCreated, gin.H{
		"success": true,
		"data":    response,
		"message": "Student registered successfully",
	})
}

// Login - Login a user
func Login(c *gin.Context) {
	var req models.LoginRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		helpers.Respond(c, false, nil, err.Error())
		return
	}

	// Find user by username with role
	var user models.User
	if err := database.DB.Preload("Role").Where("user_name = ?", req.UserName).First(&user).Error; err != nil {
		helpers.Respond(c, false, nil, "Invalid username or password")
		return
	}

	// Check if user is active
	if !user.IsActive {
		helpers.Respond(c, false, nil, "Account is deactivated")
		return
	}

	// Check password
	if !helpers.CheckPassword(user.PasswordHash, req.Password) {
		helpers.Respond(c, false, nil, "Invalid username or password")
		return
	}

	// Get role value
	roleValue := ""
	if user.Role != nil {
		roleValue = user.Role.RoleValue
	}

	// Check if user is a teacher and get teacher_id
	var teacherID *uint
	if roleValue == "teacher" {
		var teacher models.Teacher
		if err := database.DB.Where("user_id = ?", user.ID).First(&teacher).Error; err == nil {
			teacherID = &teacher.ID
		}
	}

	// Generate token with role_value and teacher_id
	token, err := helpers.GenerateToken(user.ID, user.UserName, user.RoleID, roleValue, teacherID)
	if err != nil {
		helpers.Respond(c, false, nil, "Failed to generate token")
		return
	}

	response := models.AuthResponse{
		Token: token,
		User:  user,
	}

	helpers.Respond(c, true, response, "Login successful")
}

// GetMe - Get current authenticated user
func GetMe(c *gin.Context) {
	// Get user from context (set by auth middleware)
	userID, exists := c.Get("user_id")
	if !exists {
		helpers.Respond(c, false, nil, "Unauthorized")
		return
	}

	var user models.User
	if err := database.DB.Preload("Role").First(&user, userID).Error; err != nil {
		helpers.Respond(c, false, nil, "User not found")
		return
	}

	helpers.Respond(c, true, user, "User retrieved successfully")
}

func GenerateUniqueUserName(c *gin.Context) {
	firstName := c.Query("first_name")
	lastName := c.Query("last_name")

	if lastName == "" || firstName == "" {
		helpers.Respond(c, false, nil, "First name and last name are required")
		return
	}

	baseUsername := strings.ToLower(firstName + lastName)
	baseUsername = strings.ReplaceAll(baseUsername, " ", "")

	var user models.User
	username := baseUsername
	counter := 1

	for {
		err := database.DB.Where("user_name = ?", username).First(&user).Error
		if err != nil {
			break
		}
		username = baseUsername + strconv.Itoa(counter)
		counter++
	}

	helpers.Respond(c, true, gin.H{"username": username}, "Username generated successfully") // ← Fixed
}
func Logout(c *gin.Context) {
	helpers.Respond(c, true, nil, "Logged out successfully")
}

func ChangePassword(c *gin.Context) {
	var user models.User
	id, _ := strconv.Atoi(c.Param("id"))
	var newPassword = c.PostForm("new_password")

	if err := database.DB.First(&user, id).Error; err != nil {
		helpers.RespondNotFound(c, "User not found")
		return
	}

	hashedPassword, err := helpers.HashPassword(newPassword)
	if err != nil {
		helpers.RespondInternalError(c, false, "Failed to update user password")
		return
	}

	if err := database.DB.Model(&user).Update("password_hash", hashedPassword).Error; err != nil {
		helpers.RespondBadRequest(c, "Failed to update user password")
		return
	}
	helpers.RespondCreated(c, user, "Passowrd Updated successfully")

}
