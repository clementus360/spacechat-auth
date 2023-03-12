package models

import "gorm.io/gorm"

type User struct {
	gorm.Model
	Name         string `json:"name"`
	Phone        string `json:"phone" gorm:"type:varchar(50);unique"`
	Email        string `json:"email" gorm:"type:varchar(100)"`
	Activated    bool
	TotpSecret		 string
	PhoneHash	string `json:"-" gorm:"uniqueIndex"`
}
