package controllers

import (
	"ecommerce-api/config"
	"ecommerce-api/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

func GetProductCategories(c *gin.Context) {
	rows, err := config.DB.Query(` 
		SELECT id, name, description, is_active, created_at, updated_at
		FROM product_categories
		ORDER BY id DESC
	`)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengambil data"})
		return
	}
	defer rows.Close()

	var product_categories []models.ProductCategory
	for rows.Next() {
		var category models.ProductCategory
		rows.Scan(
			&category.ID,
			&category.Name,
			&category.Description,
			&category.IsActive,
			&category.CreatedAt,
			&category.UpdatedAt,
		)
		product_categories = append(product_categories, category)
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "List Category",
		"data":    product_categories,
	})
}

func CreateProductCategory(c *gin.Context) {
	var req models.ProductCategory

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON"})
		return
	}

	if req.Name == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Nama produk category wajib diisi"})
		return
	}

	_, err := config.DB.Exec(`
		INSERT INTO product_categories (name, description, is_active, created_at)
		VALUES ($1, $2, $3, NOW())
	`, req.Name, req.Description, true)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Category sudah ada"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Category berhasil dibuat"})
}

func GetProductCategoryByID(c *gin.Context) {
	id := c.Param("id")

	var category models.ProductCategory
	err := config.DB.QueryRow(`
		SELECT id, name, description, is_active, created_at, updated_at
		FROM product_categories WHERE id = $1
	`, id).Scan(
		&category.ID,
		&category.Name,
		&category.Description,
		&category.IsActive,
		&category.CreatedAt,
		&category.UpdatedAt,
	)

	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Category tidak ditemukan"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Detail Product Category",
		"data":    category,
	})
}

func UpdateProductCategory(c *gin.Context) {
	id := c.Param("id")

	var req models.ProductCategory
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON"})
		return
	}

	_, err := config.DB.Exec(`
		UPDATE product_categories
		SET name=$1, description=$2, is_active=$3, updated_at=NOW()
		WHERE id=$4
	`, req.Name, req.Description, req.IsActive, id)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Gagal update category"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Product Category berhasil diupdate"})
}

func DeleteProductCategory(c *gin.Context) {
	id := c.Param("id")

	_, err := config.DB.Exec(`DELETE FROM product_categories WHERE id=$1`, id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Gagal hapus product category"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Product Category berhasil dihapus"})
}
