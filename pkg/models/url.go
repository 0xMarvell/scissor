package models

import "gorm.io/gorm"

type URL struct {
	gorm.Model
	originalURL  string `json:"original_url"`
	shortenedURL string `json:"shortened_url" gorm:"unique"`
	UserID       uint
}
