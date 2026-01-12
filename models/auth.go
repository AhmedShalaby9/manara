package models

type LoginRequest struct {
	UserName string `json:"user_name" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// RegisterStudentRequest - For mobile app student registration
type RegisterStudentRequest struct {
	FirstName       string `json:"first_name" binding:"required"`
	LastName        string `json:"last_name" binding:"required"`
	Email           string `json:"email" binding:"required,email"`
	Phone           string `json:"phone"`
	UserName        string `json:"user_name" binding:"required"`
	Password        string `json:"password" binding:"required,min=6"`
	TeacherUserName string `json:"teacher_user_name" binding:"required"` // Teacher's username to link student
	ParentPhone     string `json:"parent_phone"`
	AcademicYearID  uint   `json:"academic_year_id" binding:"required"`
}

// CreateUserRequest - For admin to create any user
type CreateUserRequest struct {
	FirstName string `json:"first_name" binding:"required"`
	LastName  string `json:"last_name" binding:"required"`
	Email     string `json:"email" binding:"required,email"`
	Phone     string `json:"phone"`
	UserName  string `json:"user_name" binding:"required"`
	Password  string `json:"password" binding:"required,min=6"`
	RoleID    uint   `json:"role_id" binding:"required"`
}

type AuthResponse struct {
	Token string `json:"token"`
	User  User   `json:"user"`
}
