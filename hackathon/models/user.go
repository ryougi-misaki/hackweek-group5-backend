package models

import (
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Name        string `gorm:"type:varchar(20);not null;collate:utf8mb4_unicode_ci"`
	Telephone   string `gorm:"varchar(30);not null;unique"`
	Password    string `gorm:"size:255;not null"`
	Icon        string `gorm:"varchar(50);collate:utf8mb4_unicode_ci"`
	Description string `gorm:"varchar(200);collate:utf8mb4_unicode_ci"`
	Gender      string `gorm:"varchar(2);collate:utf8mb4_unicode_ci"`
	Birth       string `gorm:"varchar(10)"`
	Role        int    `gorm:"type:int"` //默认为0，为1时是admin（后期风控补充）
}

type UserDto struct {
	Id          uint   `json:"id"`
	Name        string `json:"name"`
	Telephone   string `json:"telephone"`
	Icon        string `json:"icon"`
	Description string `json:"description"`
	Gender      string `json:"gender"`
	Birth       string `json:"birth"`
	Role        int    `json:"role"`
}

func ToUserDto(user User) UserDto {
	return UserDto{
		Id:          user.ID,
		Name:        user.Name,
		Telephone:   user.Telephone,
		Icon:        user.Icon,
		Description: user.Description,
		Gender:      user.Gender,
		Birth:       user.Birth,
		Role:        user.Role,
	}
}
