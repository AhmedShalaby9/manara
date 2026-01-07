package models

import (
	"time"
)

type Chapter struct {
	ID          uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	CourseID    uint      `gorm:"not null" json:"course_id"`
	TeacherID   uint      `gorm:"not null" json:"teacher_id"`
	Name        string    `gorm:"type:varchar(100);not null" json:"name"`
	Order       int       `gorm:"not null;default:1" json:"order"`
	Description string    `gorm:"type:text" json:"description"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`

	Course  *Course  `gorm:"foreignKey:CourseID;references:ID" json:"course,omitempty"`
	Teacher *Teacher `gorm:"foreignKey:TeacherID;references:ID" json:"teacher,omitempty"`
	Lessons []Lesson `gorm:"foreignKey:ChapterID" json:"lessons,omitempty"`
}

func (Chapter) TableName() string {
	return "chapters"
}

type CreateChapterRequest struct {
	CourseID    uint   `json:"course_id"`  // Optional for teachers (auto-filled from teacher record), required for admins
	TeacherID   uint   `json:"teacher_id"` // Optional for teachers (auto-filled from token), required for admins
	Name        string `json:"name" binding:"required"`
	Order       int    `json:"order"`
	Description string `json:"description"`
}

type UpdateChapterRequest struct {
	Name        string `json:"name"`
	Order       int    `json:"order"`
	Description string `json:"description"`
}
