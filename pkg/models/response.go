package models

// Response bodies for auth and user controllers
type SignupResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
	User    User   `json:"user"`
}

type LoginResponse struct {
	Success bool   `json:"success" example:"true"`
	Message string `json:"message" example:"login successful"`
}

type LogoutResponse struct {
	Success bool   `json:"success" example:"true"`
	Message string `json:"message" example:" logout successful"`
}

type GetUsersResponse struct {
	Success bool   `json:"success"`
	Count   int64  `json:"count"`
	Users   []User `json:"users"`
}

type GetUserResponse struct {
	Success bool `json:"success"`
	User    any  `json:"user"`
}

// Response bodies for url shortener controllers
type ShortenResponse struct {
	Message  string `json:"message" example:"short url created successfully"`
	ShortURL string `json:"short_url" example:"http://localhost:8080/api/v1/8KKoJCy"`
}
type GetURLsResponse struct {
	Success bool  `json:"success"`
	Count   int64 `json:"count"`
	URLs    []URL `json:"urls"`
}
