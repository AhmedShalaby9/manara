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
	search := c.Query("search")
	roleID, _ := strconv.Atoi(c.Query("role_id"))
	params := helpers.GetPaginationParams(c)
	query := database.DB.Model(&models.User{}).Preload("Teacher").Preload("Role").Preload("Student").Order("id DESC")

	if search != "" {
		searchPattern := "%" + search + "%"
		query = query.Where(
			"first_name LIKE ? OR last_name LIKE ? OR user_name LIKE ?",
			searchPattern, searchPattern, searchPattern,
		)
	}

	if roleID > 0 {
		query = query.Where("role_id = ?", roleID)
	}
	var users []models.User
	pagination, err := helpers.Paginate(query, params, &users)
	if err != nil {
		helpers.Respond(c, false, nil, "Failed to retrieve users")
		return
	}

	helpers.RespondWithPagin(c, true, users, "Users retrieved successfully", pagination)
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

func DeleteUser(c *gin.Context) {
	var user models.User
	id, _ := strconv.Atoi(c.Param("id"))

	if err := database.DB.Preload("Role").First(&user, id).Error; err != nil {
		helpers.RespondNotFound(c, "User not found")
		return
	}
	roleValueAny, _ := c.Get("role_value")
	roleValue, _ := roleValueAny.(string)
	if user.Role != nil &&
		models.IsSuperAdmin(user.Role.RoleValue) &&
		models.IsAdmin(roleValue) {
		helpers.RespondForbiden(c, "Admins cannot operate super admins")
		return
	}

	if err := database.DB.Delete(&user).Error; err != nil {
		helpers.Respond(c, false, nil, "Failed to delete user")
		return
	}

	helpers.RespondSuccess(c, nil, "User deleted successfully")
}
