package controllers

import (
	"ecommerce-api/config"
	"ecommerce-api/requests"
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
func UpdateUser(c *gin.Context) {
	id := c.Param("id")

	var req requests.UserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid JSON",
		})
		return
	}

	if req.Name == "" || req.Username == "" || req.RoleID == 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "name, username, dan role_id wajib diisi",
		})
		return
	}

	// update user
	res, err := config.DB.Exec(`
		UPDATE users
		SET
			name = $1,
			username = $2,
			email = $3,
			is_active = COALESCE($4, is_active),
			role_id = $5,
			updated_at = NOW()
		WHERE id = $6
	`,
		req.Name,
		req.Username,
		req.Email,
		req.IsActive,
		req.RoleID,
		id,
	)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Gagal update user",
		})
		return
	}

	affected, _ := res.RowsAffected()
	if affected == 0 {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "User tidak ditemukan",
		})
		return
	}

	var result responses.UserDetailResponse
	err = config.DB.QueryRow(`
		SELECT
			u.id,
			u.name,
			u.username,
			u.email,
			u.is_active,
			r.id,
			r.name,
			u.created_at,
			u.updated_at
		FROM users u
		JOIN roles r ON r.id = u.role_id
		WHERE u.id = $1
	`, id).Scan(
		&result.ID,
		&result.Name,
		&result.Username,
		&result.Email,
		&result.IsActive,
		&result.Role.ID,
		&result.Role.Name,
		&result.CreatedAt,
		&result.UpdatedAt,
	)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Gagal mengambil data user",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "User berhasil diupdate",
		"data":    result,
	})
}
