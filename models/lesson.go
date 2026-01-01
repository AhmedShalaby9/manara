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

	Chapter *Chapter       `gorm:"foreignKey:ChapterID;references:ID" json:"chapter,omitempty"`
	Teacher *Teacher       `gorm:"foreignKey:TeacherID;references:ID" json:"teacher,omitempty"`
	Files   []LessonFile   `gorm:"foreignKey:LessonID" json:"files,omitempty"`
	Videos  []LessonVideo  `gorm:"foreignKey:LessonID" json:"videos,omitempty"`
}

type LessonFile struct {
	ID        uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	LessonID  uint      `gorm:"not null;index" json:"lesson_id"`
	FileURL   string    `gorm:"type:varchar(500);not null" json:"file_url"`
	FileName  string    `gorm:"type:varchar(255);not null" json:"file_name"`
	FileType  string    `gorm:"type:varchar(50)" json:"file_type"`
	FileSize  int64     `gorm:"not null" json:"file_size"`
	Order     int       `gorm:"not null;default:1" json:"order"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`

	Lesson *Lesson `gorm:"foreignKey:LessonID;references:ID" json:"lesson,omitempty"`
}

func (LessonFile) TableName() string {
	return "lesson_files"
}

type LessonVideo struct {
	ID        uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	LessonID  uint      `gorm:"not null;index" json:"lesson_id"`
	VideoURL  string    `gorm:"type:varchar(500);not null" json:"video_url"`
	VideoName string    `gorm:"type:varchar(255);not null" json:"video_name"`
	FileSize  int64     `gorm:"not null" json:"file_size"`
	Order     int       `gorm:"not null;default:1" json:"order"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`

	Lesson *Lesson `gorm:"foreignKey:LessonID;references:ID" json:"lesson,omitempty"`
}

func (LessonVideo) TableName() string {
	return "lesson_videos"
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
