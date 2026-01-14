package models

import "time"

type User struct {
	ID       uint   `json:"id"`
	RoleID   int    `json:"role_id"`
	Username string `json:"username"`
	Password string `json:"password"`
	Name     string `json:"name"`
	Email    string `json:"email"`
	IsActive bool   `json:"is_active"`

	CreatedAt time.Time `json:"created_at"`
	CreatedBy *int      `json:"created_by"`

	UpdatedAt *time.Time `json:"updated_at"`
	UpdatedBy *int       `json:"updated_by"`
}
