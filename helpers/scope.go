package helpers

import (
	"strconv"

	"github.com/gin-gonic/gin"
)

// GetEffectiveTeacherID returns the teacher_id to use for filtering data.
// - For admins/super_admins: returns query param if provided, nil otherwise (all data)
// - For teachers: always returns their own teacher_id from token (ignores query param)
// - For students: returns their linked teacher_id from token (sees teacher's content)
// - For other roles: returns nil
func GetEffectiveTeacherID(c *gin.Context) *uint {
	roleValue, _ := c.Get("role_value")
	roleStr, _ := roleValue.(string)

	// Admins can filter by query param or see all data
	if roleStr == "admin" || roleStr == "super_admin" {
		if tid := c.Query("teacher_id"); tid != "" {
			id, err := strconv.ParseUint(tid, 10, 32)
			if err == nil {
				uid := uint(id)
				return &uid
			}
		}
		return nil // No filter = all data for admins
	}

	// Teachers and Students are scoped to teacher's data
	if roleStr == "teacher" || roleStr == "student" {
		if teacherID, exists := c.Get("teacher_id"); exists {
			tid := teacherID.(uint)
			return &tid
		}
	}

	return nil
}

// GetTeacherIDForCreate returns the teacher_id to use when creating resources.
// - For teachers: returns their own teacher_id from token
// - For admins: returns the teacher_id from request body (passed as parameter)
// Returns the effective teacher_id and whether it was successfully determined
func GetTeacherIDForCreate(c *gin.Context, requestTeacherID uint) (uint, bool) {
	roleValue, _ := c.Get("role_value")
	roleStr, _ := roleValue.(string)

	// Teachers always use their own ID
	if roleStr == "teacher" {
		if teacherID, exists := c.Get("teacher_id"); exists {
			return teacherID.(uint), true
		}
		return 0, false // Teacher without teacher_id in token - shouldn't happen
	}

	// Admins use the ID from request
	if roleStr == "admin" || roleStr == "super_admin" {
		if requestTeacherID > 0 {
			return requestTeacherID, true
		}
		return 0, false // Admin must provide teacher_id
	}

	return 0, false
}

// IsAdminOrSuperAdmin checks if the current user is an admin or super_admin
func IsAdminOrSuperAdmin(c *gin.Context) bool {
	roleValue, _ := c.Get("role_value")
	roleStr, _ := roleValue.(string)
	return roleStr == "admin" || roleStr == "super_admin"
}

// IsTeacher checks if the current user is a teacher
func IsTeacher(c *gin.Context) bool {
	roleValue, _ := c.Get("role_value")
	roleStr, _ := roleValue.(string)
	return roleStr == "teacher"
}
