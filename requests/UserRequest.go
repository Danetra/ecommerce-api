package requests

type UserRequest struct {
	Name     string  `json:"name"`
	Username *string `json:"username"`
	Email    *string `json:"email"`
	IsActive *bool   `json:"is_active"`
	RoleID   int     `json:"role_id"`
}
