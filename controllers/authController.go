package controllers

import (
	"ecommerce-api/config"
	helper "ecommerce-api/helpers"
	"ecommerce-api/models"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

// Login godoc
// @Summary Login user
// @Description Login dengan username & password
// @Tags Auth
// @Accept json
// @Produce json
// @Param body body models.LoginRequest true "Login Request"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /api/login [post]
func Login(c *gin.Context) {
	var req struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON"})
		return
	}

	var user models.User
	err := config.DB.QueryRow(
		"SELECT id,username,password FROM users WHERE username=$1 AND is_active = true",
		req.Username,
	).Scan(&user.ID, &user.Username, &user.Password)

	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Username atau password salah"})
		return
	}

	if bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)) != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Username atau password salah"})
		return
	}

	token, expiredAt, err := helper.GenerateJWT(int(user.ID), user.Username)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Gagal generate token",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Login berhasil",
		"data": gin.H{
			"type":       "Bearer",
			"token":      token,
			"expired_in": 3600,
			"expired_at": expiredAt.Format(time.RFC3339),
		},
	})
}

// Register godoc
// @Summary Register user
// @Description Register user baru
// @Tags Auth
// @Accept json
// @Produce json
// @Param body body models.RegisterRequest true "Register Request"
// @Success 201 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /api/register [post]
func Register(c *gin.Context) {
	var req struct {
		RoleID   int    `json:"role_id"`
		Username string `json:"username"`
		Password string `json:"password"`
		Name     string `json:"name"`
		Email    string `json:"email"`
		IsActive bool   `json:"is_active"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON"})
		return
	}

	if req.Username == "" || req.Password == "" || req.Name == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "username, password, dan name wajib diisi",
		})
		return
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Gagal memproses password",
		})
		return
	}

	if !req.IsActive {
		req.IsActive = true
	}

	var userID int

	err = config.DB.QueryRow(`
		INSERT INTO users
			(role_id, username, password, name, email, is_active, created_at)
		VALUES
			($1, $2, $3, $4, $5, $6, NOW())
		RETURNING id
	`,
		req.RoleID,
		req.Username,
		string(hash),
		req.Name,
		req.Email,
		req.IsActive,
	).Scan(&userID)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Username atau email sudah terdaftar",
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "User Register Success",
		"data": gin.H{
			"id":       userID,
			"username": req.Username,
			"name":     req.Name,
		},
	})
}
