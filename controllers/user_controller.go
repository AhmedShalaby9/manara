package controllers

import (
	"manara/database"
	"manara/helpers"
	"manara/models"
	"strconv"

	"github.com/gin-gonic/gin"
)

func CreateUser(c *gin.Context) {
	var req models.RegisterRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		helpers.Respond(c, false, nil, err.Error())
		return
	}

	var registeredUser models.User

	res := database.DB.Where("user_name =?", req.UserName).First(&registeredUser)
	if res.Error == nil {
		helpers.Respond(c, false, nil, "User already exists")
		return
	}

	if err := database.DB.Where("email =?", req.Email).First(&registeredUser).Error; err == nil {
		helpers.Respond(c, false, nil, "Email already exists")
		return

	}
	hashedPassword, err := helpers.HashPassword(req.Password)
	if err != nil {
		helpers.Respond(c, false, nil, "Failed to hash password")
		return
	}
	user := models.User{
		FirstName:    req.FirstName,
		LastName:     req.LastName,
		Email:        req.Email,
		Phone:        req.Phone,
		UserName:     req.UserName,
		PasswordHash: hashedPassword,
		RoleID:       req.RoleID,
		IsActive:     false,
	}

	if err := database.DB.Create(&user).Error; err != nil {
		helpers.Respond(c, false, nil, "Failed to create user")
		return
	}

	database.DB.Preload("Role").First(&user, user.ID)
	helpers.Respond(c, true, user, "User created successfully")

}

func GetUsers(c *gin.Context) {
	var users []models.User
	if err := database.DB.Preload("Role").Find(&users).Error; err != nil {
		helpers.Respond(c, false, nil, err.Error())
		return
	}
	helpers.Respond(c, true, users, "Failed to retrive users")

}
func ToggleUserActivation(c *gin.Context) {
	userID := c.Param("id")
	activeParam := c.PostForm("active")
	isActive, err := strconv.ParseBool(activeParam)

	if err != nil {
		helpers.Respond(c, false, nil, "Invalid active value")
		return
	}

	var user models.User

	if err := database.DB.First(&user, userID).Error; err != nil {
		helpers.Respond(c, false, nil, "User not found")
		return
	}

	if err := database.DB.Model(&user).
		Update("is_active", isActive).Error; err != nil {
		helpers.Respond(c, false, nil, "Failed to update user")
		return
	}
	database.DB.Preload("Role").First(&user)

	helpers.Respond(c, true, user, "User updated successfully")
}
