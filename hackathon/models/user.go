package models

import "gorm.io/gorm"

type User struct {
	gorm.Model
	Name        string `gorm:"type:varchar(20);not null"`
	Telephone   string `gorm:"varchar(30);not null;unique"`
	Password    string `gorm:"size:255;not null"`
	Icon        string `gorm:"varchar(50);"`
	Description string `gorm:"varchar(200)"`
}
