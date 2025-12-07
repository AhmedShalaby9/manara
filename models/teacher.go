package models

import (
	"time"

	"gorm.io/gorm"
)

type Teacher struct {
	ID              uint           `gorm:"primaryKey;autoIncrement" json:"id"`
	UserID          uint           `gorm:"unique;not null" json:"user_id"`
	Bio             string         `gorm:"type:text" json:"bio"`
	Specialization  string         `gorm:"type:varchar(255)" json:"specialization"`
	ExperienceYears int            `gorm:"default:0" json:"experience_years"`
	CreatedAt       time.Time      `json:"created_at"`
	UpdatedAt       time.Time      `json:"updated_at"`
	DeletedAt       gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`

	User     User      `gorm:"foreignKey:UserID;references:ID" json:"user,omitempty"`
	Students []Student `gorm:"foreignKey:TeacherID" json:"students,omitempty"`
}

func (Teacher) TableName() string {
	return "teachers"
}
