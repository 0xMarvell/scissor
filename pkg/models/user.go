package models

import "gorm.io/gorm"

type User struct {
	gorm.Model
	Username string `json:"user_name" gorm:"unique"`
	Password string `json:"-"`
	URLs     []URL  `gorm:"foreignKey:UserID" json:"urls"`
}
