package controllers

import (
	"ecommerce-api/config"
	"ecommerce-api/responses"
	"net/http"

	"github.com/gin-gonic/gin"
)

func GetUserByID(c *gin.Context) {
	id := c.Param("id")

	var user responses.UserDetailResponse
	err := config.DB.QueryRow(`
		SELECT 
			u.id,
			u.name,
			u.username,
			u.email,
			u.is_active,
			u.role_id,
			r.name AS role_name,
			u.created_at,
			u.updated_at
		FROM users u
		JOIN roles r ON r.id = u.role_id
		WHERE u.id = $1
	`, id).Scan(
		&user.ID,
		&user.Name,
		&user.Username,
		&user.Email,
		&user.IsActive,
		&user.Role.ID,
		&user.Role.Name,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "User tidak ditemukan",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Detail User",
		"data":    user,
	})
}
