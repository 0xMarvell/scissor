package models

// Request bodies for auth and user controllers
type SignupRequest struct {
	Username string `json:"username" example:"linchung"`
	Password string `json:"password" example:"hero108"`
}
type LoginRequest struct {
	Username string `json:"username" example:"linchung"`
	Password string `json:"password" example:"hero108"`
}

// Request bodies for url shortener controllers
type ShortenURLRequest struct {
	URL string `json:"url" example:"https://google.com"`
}
