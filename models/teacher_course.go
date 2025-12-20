package models

import (
	"time"
)

type TeacherCourse struct {
	ID        uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	TeacherID uint      `gorm:"not null" json:"teacher_id"`
	CourseID  uint      `gorm:"not null" json:"course_id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`

	Teacher *Teacher `gorm:"foreignKey:TeacherID;references:ID" json:"teacher,omitempty"`
	Course  *Course  `gorm:"foreignKey:CourseID;references:ID" json:"course,omitempty"`
}

func (TeacherCourse) TableName() string {
	return "teacher_courses"
}

type AssignCourseToTeacherRequest struct {
	TeacherID uint `json:"teacher_id" binding:"required"`
	CourseID  uint `json:"course_id" binding:"required"`
}
