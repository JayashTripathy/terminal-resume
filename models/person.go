package models

import "gorm.io/gorm"

type URL struct {
	Label string `json:"label"`
	Href  string `json:"href"`
}


type Person struct {
	gorm.Model
	UserID  uint `json:"user_id"`
	Headline string `json:"headline"`
	Phone string `json:"phone"`
	Location string `json:"location"`
	Url URL `gorm:"embedded" json:"url"`
}