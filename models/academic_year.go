package models

import (
	"time"
)

type AcademicYear struct {
	ID        uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	Name      string    `gorm:"type:varchar(100);not null" json:"name"`
	IsActive  bool      `gorm:"type:boolean;not null;default:true" json:"is_active"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (AcademicYear) TableName() string {
	return "academicYears"
}
