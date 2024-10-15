package models

import "gorm.io/gorm"

type User struct {
	gorm.Model
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
	Role     string `gorm:"default:'user'" json:"role"`
	Person   Person `gorm:"foreignKey:UserID" json:"person"`
}

