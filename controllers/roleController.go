package controllers

import (
	"ecommerce-api/config"
	"ecommerce-api/models"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

func GetRoles(c *gin.Context) {
	rows, err := config.DB.Query(`
		SELECT id, name, description, created_at, updated_at
		FROM roles
		ORDER BY id ASC
	`)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Gagal mengambil data role",
		})
		return
	}
	defer rows.Close()

	var roles []models.Role

	for rows.Next() {
		var r models.Role
		if err := rows.Scan(
			&r.ID,
			&r.Name,
			&r.Description,
			&r.CreatedAt,
			&r.UpdatedAt,
		); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Gagal membaca data role",
			})
			return
		}
		roles = append(roles, r)
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "List roles",
		"data":    roles,
	})
}

func CreateRole(c *gin.Context) {
	var req struct {
		Name        string `json:"name"`
		Description string `json:"description"`
	}

	if err := c.ShouldBindJSON(&req); err != nil || req.Name == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Nama role wajib diisi",
		})
		return
	}

	var roleID int
	err := config.DB.QueryRow(`
		INSERT INTO roles (name, description, created_at)
		VALUES ($1, $2, NOW())
		RETURNING id
	`, req.Name, req.Description).Scan(&roleID)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Role sudah ada",
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Role berhasil dibuat",
		"data": gin.H{
			"id":   roleID,
			"name": req.Name,
		},
	})
}

func GetRoleByID(c *gin.Context) {
	id := c.Param("id")

	var role models.Role
	err := config.DB.QueryRow(`
		SELECT id, name, description, created_at, updated_at
		FROM roles
		WHERE id = $1
	`, id).Scan(
		&role.ID,
		&role.Name,
		&role.Description,
		&role.CreatedAt,
		&role.UpdatedAt,
	)

	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Role tidak ditemukan",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Detail role",
		"data":    role,
	})
}

func UpdateRole(c *gin.Context) {
	id := c.Param("id")

	var req struct {
		Name        string `json:"name"`
		Description string `json:"description"`
	}

	if err := c.ShouldBindJSON(&req); err != nil || req.Name == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Nama role wajib diisi",
		})
		return
	}

	res, err := config.DB.Exec(`
		UPDATE roles
		SET name = $1,
		    description = $2,
		    updated_at = $3
		WHERE id = $4
	`, req.Name, req.Description, time.Now(), id)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Gagal update role",
		})
		return
	}

	affected, _ := res.RowsAffected()
	if affected == 0 {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Role tidak ditemukan",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Role berhasil diupdate",
	})
}
func DeleteRole(c *gin.Context) {
	id := c.Param("id")

	res, err := config.DB.Exec(`
		DELETE FROM roles WHERE id = $1
	`, id)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Role tidak bisa dihapus (masih digunakan)",
		})
		return
	}

	affected, _ := res.RowsAffected()
	if affected == 0 {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Role tidak ditemukan",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Role berhasil dihapus",
	})
}
