package responses

type UserDetailResponse struct {
	ID        int          `json:"id"`
	Name      string       `json:"name"`
	Username  string       `json:"username"`
	Email     *string      `json:"email"`
	IsActive  bool         `json:"is_active"`
	Role      RoleResponse `json:"role"`
	CreatedAt string       `json:"created_at"`
	UpdatedAt *string      `json:"updated_at"`
}
