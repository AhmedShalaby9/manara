package models

import (
	"time"
)

type User struct {
	ID           uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	FirstName    string    `gorm:"type:varchar(100);not null" json:"first_name"`
	LastName     string    `gorm:"type:varchar(100);not null" json:"last_name"`
	RoleID       uint      `gorm:"not null" json:"role_id"`
	Email        string    `gorm:"type:varchar(255);unique;not null" json:"email"`
	Phone        string    `gorm:"type:varchar(20)" json:"phone"`
	UserName     string    `gorm:"type:varchar(100);unique;not null" json:"user_name"`
	PasswordHash string    `gorm:"type:varchar(255);not null" json:"-"`
	IsActive     bool      `gorm:"default:true" json:"is_active"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`

	Role    *Role    `gorm:"foreignKey:RoleID;references:ID" json:"role,omitempty"`
	Teacher *Teacher `gorm:"foreignKey:UserID" json:"teacher,omitempty"`
	Student *Student `gorm:"foreignKey:UserID" json:"student,omitempty"`
}

func (User) TableName() string {
	return "users"
}
