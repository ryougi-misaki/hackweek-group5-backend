package models

import (
	"gorm.io/gorm"
)

type Tag struct {
	gorm.Model
	Name        string `gorm:"type:varchar(20);not null;collate:utf8mb4_unicode_ci"`
	Description string `gorm:"type:varchar(200);collate:utf8mb4_unicode_ci"`
}

type Post struct {
	gorm.Model
	UserId  uint   `gorm:"type:bigint unsigned;not null"`
	TagId   uint   `gorm:"type:bigint unsigned;not null"`
	Title   string `gorm:"type:varchar(20);collate:utf8mb4_unicode_ci"`
	Content string `gorm:"type:longtext;collate:utf8mb4_unicode_ci"`
}
