package models

import (
	"gorm.io/gorm"
)

type ChatRecord struct {
	gorm.Model
	UserId  uint   `gorm:"type:int;not null"`
	FromId  uint   `gorm:"type:int;not null"`
	LastMsg string `gorm:"type:varchar(100);collate:utf8mb4_unicode_ci"`
	HasBg   int    `gorm:"type:int"` //是否置顶，默认0不置顶，1置顶
}
