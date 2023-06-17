package models

import "gorm.io/gorm"

// URL represents a url object
//
// swagger:model
type URL struct {
	gorm.Model
	// The url to be shortened
	//
	// required: true
	// example: https://google.com/
	OriginalURL string `json:"original_url"`
	// The shortened version of the original url
	//
	// required: false
	// unique: true
	// example: djna7L
	ShortenedURL string `json:"shortened_url" gorm:"unique"`
	// The user ID of the user that created the shortened url
	//
	// example: 1
	UserID uint `json:"user_id"`
}
