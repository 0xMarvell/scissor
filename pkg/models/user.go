package models

import (
	"gorm.io/gorm"
)

// User represents a user object
//
// swagger:model
type User struct {
	gorm.Model
	// The username of the user
	//
	// required: true
	// unique: true
	// example: john_doe
	Username string `json:"username" gorm:"unique"`
	// The password of the user
	//
	// required: true
	// writeOnly: true
	// example: password123
	Password string `json:"-"`
	// The urls created by the user
	//
	// required: false
	// example: password123
	URLs []URL `json:"urls"`
}
