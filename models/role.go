package models

import (
	"time"

	"gorm.io/gorm"
)

type Role struct {
	ID        uint           `gorm:"primaryKey;autoIncrement" json:"id"`
	RoleName  string         `gorm:"type:varchar(100);not null" json:"role_name"`
	RoleValue string         `gorm:"type:varchar(100);unique;not null" json:"role_value"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`

	Users []User `gorm:"foreignKey:RoleID" json:"users,omitempty"`
}

func (Role) TableName() string {
	return "roles"
}
