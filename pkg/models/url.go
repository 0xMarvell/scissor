package models

import "gorm.io/gorm"

type URL struct {
	gorm.Model
	OriginalURL  string `json:"original_url"`
	ShortenedURL string `json:"shortened_url" gorm:"unique"`
	UserID       uint   `json:"user_id"`
}
