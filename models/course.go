package models

import (
	"time"
)

type Course struct {
	ID          uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	Name        string    `gorm:"type:varchar(100);not null" json:"name"`
	Description string    `gorm:"type:text" json:"description"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	ImageURL    string    `gorm:"type:varchar(500)" json:"image_url"`

	Chapters []Chapter `gorm:"foreignKey:CourseID" json:"chapters,omitempty"`
	Teachers []Teacher `gorm:"foreignKey:CourseID" json:"teachers,omitempty"`
}

func (Course) TableName() string {
	return "courses"
}

type UpdateCourseRequest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}
