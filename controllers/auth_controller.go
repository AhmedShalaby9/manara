package controllers

import (
	"manara/database"
	"manara/helpers"
	"manara/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

// Register - Register a new user
func Register(c *gin.Context) {
	var req models.RegisterRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		helpers.Respond(c, false, nil, err.Error())
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

	// Create user
	user := models.User{
		FirstName:    req.FirstName,
		LastName:     req.LastName,
		Email:        req.Email,
		Phone:        req.Phone,
		UserName:     req.UserName,
		PasswordHash: hashedPassword,
		RoleID:       req.RoleID,
		IsActive:     true,
	}

	if err := database.DB.Create(&user).Error; err != nil {
		helpers.Respond(c, false, nil, "Failed to create user")
		return
	}

	// Generate token
	token, err := helpers.GenerateToken(user.ID, user.UserName, user.RoleID)
	if err != nil {
		helpers.Respond(c, false, nil, "Failed to generate token")
		return
	}

	// Load user with role
	database.DB.Preload("Role").First(&user, user.ID)

	response := models.AuthResponse{
		Token: token,
		User:  user,
	}

	c.JSON(http.StatusCreated, gin.H{
		"success": true,
		"data":    response,
		"message": "User registered successfully",
	})
}

// Login - Login a user
func Login(c *gin.Context) {
	var req models.LoginRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		helpers.Respond(c, false, nil, err.Error())
		return
	}

	// Find user by username
	var user models.User
	if err := database.DB.Where("user_name = ?", req.UserName).First(&user).Error; err != nil {
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

	// Generate token
	token, err := helpers.GenerateToken(user.ID, user.UserName, user.RoleID)
	if err != nil {
		helpers.Respond(c, false, nil, "Failed to generate token")
		return
	}

	// Load user with role
	database.DB.Preload("Role").First(&user, user.ID)

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

func Logout(c *gin.Context) {
	helpers.Respond(c, true, nil, "Logged out successfully")
}
