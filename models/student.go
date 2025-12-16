package models

import (
	"time"
)

type Student struct {
	ID             uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	UserID         uint      `gorm:"unique;not null" json:"user_id"`
	TeacherID      uint      `gorm:"not null" json:"teacher_id"`
	AcademicYearID uint      `gorm:"not null;default:3" json:"academic_year_id"` // ← ADD THIS
	ParentPhone    string    `gorm:"type:varchar(20)" json:"parent_phone"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`

	User         *User         `gorm:"foreignKey:UserID;references:ID" json:"user,omitempty"`
	Teacher      *Teacher      `gorm:"foreignKey:TeacherID;references:ID" json:"teacher,omitempty"`
	AcademicYear *AcademicYear `gorm:"foreignKey:AcademicYearID;references:ID" json:"academic_year,omitempty"` // ← ADD THIS
}

func (Student) TableName() string {
	return "students"
}

type CreateStudentRequest struct {
	FirstName      string `json:"first_name" binding:"required"`
	LastName       string `json:"last_name" binding:"required"`
	Email          string `json:"email" binding:"required,email"`
	Phone          string `json:"phone"`
	UserName       string `json:"user_name" binding:"required"`
	Password       string `json:"password" binding:"required,min=6"`
	TeacherID      uint   `json:"teacher_id" binding:"required"`
	AcademicYearID uint   `json:"academic_year_id" binding:"required"`
	ParentPhone    string `json:"parent_phone"`
}
