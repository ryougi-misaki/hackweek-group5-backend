package models

import "gorm.io/gorm"

type User struct {
	gorm.Model
	Name        string `gorm:"type:varchar(20);not null"`
	Telephone   string `gorm:"varchar(30);not null;unique"`
	Password    string `gorm:"size:255;not null"`
	Icon        string `gorm:"varchar(50);"`
	Description string `gorm:"varchar(200)"`
	Gender      string `gorm:"varchar(2)"`
	Birth       string `gorm:"varchar(10)"`
}

type UserDto struct {
	Id        uint   `json:"id"`
	Name      string `json:"name"`
	Telephone string `json:"telephone"`
}

func ToUserDto(user User) UserDto {
	return UserDto{
		Id:        user.ID,
		Name:      user.Name,
		Telephone: user.Telephone,
	}
}
