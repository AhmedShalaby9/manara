package models

import (
	"time"
)

type Teacher struct {
	ID              uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	UserID          uint      `gorm:"unique;not null" json:"user_id"`
	Bio             string    `gorm:"type:text" json:"bio"`
	Specialization  string    `gorm:"type:varchar(255)" json:"specialization"`
	ExperienceYears int       `gorm:"default:0" json:"experience_years"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`

	User     User      `gorm:"foreignKey:UserID;references:ID" json:"user,omitempty"`
	Students []Student `gorm:"foreignKey:TeacherID" json:"students,omitempty"`
}

func (Teacher) TableName() string {
	return "teachers"
}

type CreateTeacherRequest struct {
	FirstName       string `json:"first_name" binding:"required"`
	LastName        string `json:"last_name" binding:"required"`
	Email           string `json:"email" binding:"required,email"`
	Phone           string `json:"phone"`
	UserName        string `json:"user_name" binding:"required"`
	Password        string `json:"password" binding:"required,min=6"`
	Bio             string `json:"bio"`
	Specialization  string `json:"specialization"`
	ExperienceYears int    `json:"experience_years"`
}

type UpdateTeacherRequest struct {
	Bio             string `json:"bio"`
	Specialization  string `json:"specialization"`
	ExperienceYears int    `json:"experience_years"`
}
