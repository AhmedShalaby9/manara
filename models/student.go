package models

import (
	"time"
)

type Student struct {
	ID          uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	UserID      uint      `gorm:"unique;not null" json:"user_id"`
	TeacherID   uint      `gorm:"not null" json:"teacher_id"`
	GradeLevel  uint      `gorm:"type:varchar(3)" json:"grade_level"`
	ParentPhone string    `gorm:"type:varchar(20)" json:"parent_phone"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`

	User    *User    `gorm:"foreignKey:UserID;references:ID" json:"user,omitempty"`
	Teacher *Teacher `gorm:"foreignKey:TeacherID;references:ID" json:"teacher,omitempty"`
}

func (Student) TableName() string {
	return "students"
}

type CreateStudentRequest struct {
	FirstName   string `json:"first_name" binding:"required"`
	LastName    string `json:"last_name" binding:"required"`
	Email       string `json:"email" binding:"required,email"`
	Phone       string `json:"phone"`
	UserName    string `json:"user_name" binding:"required"`
	Password    string `json:"password" binding:"required,min=6"`
	GradeLevel  uint   `json:"grade_level"`
	ParentPhone string `json:"parent_phone"`
	TeacherId   uint   `json:"teacher_id"`
}

// type UpdateStudentRequest struct {
// 	gradeLevel  string `json:"grade_level"`
// 	parentPhone string `json:"parent_phone"`
// }
