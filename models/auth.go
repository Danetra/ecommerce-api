package models

type LoginRequest struct {
	Username string `json:"username" example:"administrator"`
	Password string `json:"password" example:"admin123"`
}

type RegisterRequest struct {
	RoleID   int    `json:"role_id" example:"1"`
	Username string `json:"username" example:"john_doe"`
	Password string `json:"password" example:"password123"`
	Name     string `json:"name" example:"John Doe"`
	Email    string `json:"email" example:"john@mail.com"`
	IsActive bool   `json:"is_active" example:"true"`
}
