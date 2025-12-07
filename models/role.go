package models

import (
	"time"
)

type Role struct {
	ID        uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	RoleName  string    `gorm:"type:varchar(100);not null" json:"role_name"`
	RoleValue string    `gorm:"type:varchar(100);unique;not null" json:"role_value"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`

	Users []User `gorm:"foreignKey:RoleID" json:"users,omitempty"`
}

func (Role) TableName() string {
	return "roles"
}
