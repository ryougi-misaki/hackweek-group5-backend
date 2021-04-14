package models

import (
	"gorm.io/gorm"
)

type Tag struct {
	gorm.Model
	Name        string `gorm:"type:varchar(20);not null"`
	Description string `gorm:"type:varchar(200)"`
}

type Post struct {
	gorm.Model
	UserId  uint   `gorm:"type:bigint unsigned;not null"`
	TagId   uint   `gorm:"type:bigint unsigned;not null"`
	Title   string `gorm:"type:varchar(20)"`
	Content string `gorm:"type:longtext"`
}
