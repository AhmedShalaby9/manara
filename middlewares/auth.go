package middleware

import (
	"manara/database"
	"manara/helpers"
	"manara/models"
	"strings"

	"github.com/gin-gonic/gin"
)

// AuthMiddleware validates JWT token
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			helpers.RespondUnauthorized(c, "Unauthorized")
			c.Abort()
			return
		}

		// Extract token from "Bearer <token>"
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			helpers.Respond(c, false, nil, "Invalid authorization header format")
			c.Abort()
			return
		}

		tokenString := parts[1]

		// Validate token
		claims, err := helpers.ValidateToken(tokenString)
		if err != nil {
			helpers.RespondTokenError(c, "Invalid or expired token")
			c.Abort()
			return
		}

		// Set user info in context
		c.Set("user_id", claims.UserID)
		c.Set("user_name", claims.UserName)
		c.Set("role_id", claims.RoleID)
		c.Set("role_value", claims.RoleValue)
		if claims.TeacherID != nil {
			c.Set("teacher_id", *claims.TeacherID)
		}

		c.Next()
	}
}

// RoleMiddleware checks if user has one of the allowed roles
func RoleMiddleware(allowedRoleValues ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		roleID, exists := c.Get("role_id")
		if !exists {
			helpers.Respond(c, false, nil, "Unauthorized")
			c.Abort()
			return
		}

		// Get role from database
		var role models.Role
		if err := database.DB.First(&role, roleID).Error; err != nil {
			helpers.Respond(c, false, nil, "Invalid role")
			c.Abort()
			return
		}

		// Check if user's role is in allowed roles
		allowed := false
		for _, allowedRole := range allowedRoleValues {
			if role.RoleValue == allowedRole {
				allowed = true
				break
			}
		}

		if !allowed {
			helpers.RespondForbiden(c, "Insufficient permissions")
			c.Abort()
			return
		}

		c.Next()
	}
}
