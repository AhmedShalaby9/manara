package models

import (
	"time"

	"gorm.io/gorm"
)

type Student struct {
	ID          uint           `gorm:"primaryKey;autoIncrement" json:"id"`
	UserID      uint           `gorm:"unique;not null" json:"user_id"`
	TeacherID   uint           `gorm:"not null" json:"teacher_id"`
	GradeLevel  string         `gorm:"type:varchar(50)" json:"grade_level"`
	ParentPhone string         `gorm:"type:varchar(20)" json:"parent_phone"`
	DateOfBirth *time.Time     `gorm:"type:date" json:"date_of_birth,omitempty"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`

	User    User    `gorm:"foreignKey:UserID;references:ID" json:"user,omitempty"`
	Teacher Teacher `gorm:"foreignKey:TeacherID;references:ID" json:"teacher,omitempty"`
}

func (Student) TableName() string {
	return "students"
}
