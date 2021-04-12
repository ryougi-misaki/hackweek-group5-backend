package models

import (
	"gorm.io/gorm"
)

type SmsCode struct {
	gorm.Model
	Phone string `gorm:"varchar(11)" json:"phone"`
	BizId string `gorm:"varchar(30)" json:"biz_id"`
	Code  string `gorm:"varchar(6)" json:"code"`
}
