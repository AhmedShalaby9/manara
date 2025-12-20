package models

import (
	"time"
)

type Lesson struct {
	ID          uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	ChapterID   uint      `gorm:"not null" json:"chapter_id"`
	TeacherID   uint      `gorm:"not null" json:"teacher_id"`
	Name        string    `gorm:"type:varchar(200);not null" json:"name"`
	Description string    `gorm:"type:text" json:"description"`
	Order       int       `gorm:"not null;default:1" json:"order"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`

	Chapter *Chapter `gorm:"foreignKey:ChapterID;references:ID" json:"chapter,omitempty"`
	Teacher *Teacher `gorm:"foreignKey:TeacherID;references:ID" json:"teacher,omitempty"`
}

func (Lesson) TableName() string {
	return "lessons"
}

type CreateLessonRequest struct {
	ChapterID   uint   `json:"chapter_id" binding:"required"`
	TeacherID   uint   `json:"teacher_id" binding:"required"`
	Name        string `json:"name" binding:"required"`
	Description string `json:"description"`
	Order       int    `json:"order"`
}

type UpdateLessonRequest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Order       int    `json:"order"`
}
