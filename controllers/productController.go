package controllers

import (
	"ecommerce-api/config"
	"ecommerce-api/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

func GetProducts(c *gin.Context) {
	rows, err := config.DB.Query(`
		SELECT 
			p.id, p.name, p.description, p.price, p.stock, p.image, p.is_active, p.created_at,
			c.id, c.name
		FROM products p
		JOIN categories c ON c.id = p.category_id
		WHERE p.is_active = true
		ORDER BY p.id DESC
	`)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal ambil product"})
		return
	}
	defer rows.Close()

	var products []models.Product

	for rows.Next() {
		var p models.Product
		var cty models.ProductCategory

		rows.Scan(
			&p.ID,
			&p.Name,
			&p.Description,
			&p.Price,
			&p.Stock,
			&p.Image,
			&p.IsActive,
			&p.CreatedAt,
			&cty.ID,
			&cty.Name,
		)

		p.Category = &cty
		products = append(products, p)
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "List Products",
		"data":    products,
	})
}

func CreateProduct(c *gin.Context) {
	var req models.Product

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON"})
		return
	}

	if req.Name == "" || req.Price <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Name dan price wajib diisi"})
		return
	}

	_, err := config.DB.Exec(`
		INSERT INTO products 
		(category_id, name, description, price, stock, image, is_active, created_at)
		VALUES ($1,$2,$3,$4,$5,$6,$7,NOW())
	`,
		req.CategoryID,
		req.Name,
		req.Description,
		req.Price,
		req.Stock,
		req.Image,
		true,
	)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Gagal membuat product"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Product berhasil dibuat",
	})
}

func GetProductByID(c *gin.Context) {
	id := c.Param("id")

	var p models.Product
	var cty models.ProductCategory

	err := config.DB.QueryRow(`
		SELECT 
			p.id, p.name, p.description, p.price, p.stock, p.image, p.is_active, p.created_at,
			c.id, c.name
		FROM products p
		JOIN categories c ON c.id = p.category_id
		WHERE p.id = $1
	`, id).Scan(
		&p.ID,
		&p.Name,
		&p.Description,
		&p.Price,
		&p.Stock,
		&p.Image,
		&p.IsActive,
		&p.CreatedAt,
		&cty.ID,
		&cty.Name,
	)

	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Product tidak ditemukan"})
		return
	}

	p.Category = &cty

	c.JSON(http.StatusOK, gin.H{
		"message": "Detail Product",
		"data":    p,
	})
}

func UpdateProduct(c *gin.Context) {
	id := c.Param("id")

	var req models.Product
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON"})
		return
	}

	_, err := config.DB.Exec(`
		UPDATE products
		SET category_id=$1, name=$2, description=$3,
		    price=$4, stock=$5, image=$6, is_active=$7, updated_at=NOW()
		WHERE id=$8
	`,
		req.CategoryID,
		req.Name,
		req.Description,
		req.Price,
		req.Stock,
		req.Image,
		req.IsActive,
		id,
	)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Gagal update product"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Product berhasil diupdate"})
}
func DeleteProduct(c *gin.Context) {
	id := c.Param("id")

	_, err := config.DB.Exec(`DELETE FROM products WHERE id=$1`, id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Gagal hapus product"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Product berhasil dihapus"})
}
