package controllers

import (
	"ecommerce-api/config"
	helper "ecommerce-api/helpers"
	"ecommerce-api/models"
	"ecommerce-api/responses"
	"fmt"
	"net/http"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

func GetProducts(c *gin.Context) {
	rows, err := config.DB.Query(`SELECT p.id, p.category_id, p.name, p.description, p.price, p.stock, p.image, p.is_active, p.created_at, p.updated_at, c.id, c.name, c.description FROM products p JOIN product_categories c ON c.id = p.category_id WHERE p.is_active = true ORDER BY p.id DESC`)

	if err != nil {
		fmt.Println("SQL ERROR:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal ambil product"})
		return
	}
	defer rows.Close()

	var result []responses.ProductResponse

	for rows.Next() {
		var p models.Product
		var cat models.ProductCategory

		if err := rows.Scan(
			&p.ID,
			&p.CategoryID,
			&p.Name,
			&p.Description,
			&p.Price,
			&p.Stock,
			&p.Image,
			&p.IsActive,
			&p.CreatedAt,
			&p.UpdatedAt,
			&cat.ID,
			&cat.Name,
			&cat.Description,
		); err != nil {
			fmt.Println("SQL ERROR:", err)
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Gagal parsing data",
			})
			return
		}

		resp := responses.ProductResponse{
			ID:          p.ID,
			CategoryID:  p.CategoryID,
			Name:        p.Name,
			Description: p.Description,
			Price:       p.Price,
			Stock:       p.Stock,
			Image:       helper.FileURL(c, p.Image),
			IsActive:    p.IsActive,
			CreatedAt:   p.CreatedAt,
			UpdatedAt:   p.UpdatedAt,
			Category: responses.ProductCategoryResponse{
				ID:          cat.ID,
				Name:        cat.Name,
				Description: cat.Description,
			},
		}

		result = append(result, resp)
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "List Products",
		"total":   len(result),
		"data":    result,
	})
}

func CreateProduct(c *gin.Context) {
	name := c.PostForm("name")
	description := c.PostForm("description")
	categoryID, _ := strconv.Atoi(c.PostForm("category_id"))
	price, _ := strconv.ParseFloat(c.PostForm("price"), 64)
	stock, _ := strconv.Atoi(c.PostForm("stock"))

	createdBy := c.GetInt("user_id")

	if createdBy == 0 {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "Unauthorized",
		})
		return
	}

	if name == "" || price <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Name dan price wajib diisi",
		})
		return
	}

	// upload image
	file, err := c.FormFile("image")
	var imagePath string

	if err == nil {
		ext := strings.ToLower(filepath.Ext(file.Filename))
		if ext != ".jpg" && ext != ".jpeg" && ext != ".png" && ext != ".webp" {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Format image harus jpg, jpeg, png, atau webp",
			})
			return
		}

		// generate nama file unik
		filename := fmt.Sprintf(
			"%d_%s",
			time.Now().UnixNano(),
			file.Filename,
		)

		imagePath = "uploads/products/" + filename

		// simpan file
		if err := c.SaveUploadedFile(file, imagePath); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Gagal upload image",
			})
			return
		}
	}

	_, err = config.DB.Exec(`
		INSERT INTO products 
		(category_id, name, description, price, stock, image, is_active, created_by, created_at)
		VALUES ($1,$2,$3,$4,$5,$6,$7,$8,NOW())
	`,
		categoryID,
		name,
		description,
		price,
		stock,
		imagePath,
		true,
		createdBy,
	)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Gagal membuat product",
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Product berhasil dibuat",
		"data": gin.H{
			"name":  name,
			"image": helper.FileURL(c, imagePath),
		},
	})
}

func GetProductByID(c *gin.Context) {
	id := c.Param("id")

	row := config.DB.QueryRow(`
		SELECT 
			p.id,
			p.category_id,
			p.name,
			p.description,
			p.price,
			p.stock,
			p.image,
			p.is_active,
			p.created_at,
			p.updated_at,
			c.id,
			c.name,
			c.description
		FROM products p
		JOIN product_categories c ON c.id = p.category_id
		WHERE p.id = $1
	`, id)

	var p models.Product
	var cat models.ProductCategory

	if err := row.Scan(
		&p.ID,
		&p.CategoryID,
		&p.Name,
		&p.Description,
		&p.Price,
		&p.Stock,
		&p.Image,
		&p.IsActive,
		&p.CreatedAt,
		&p.UpdatedAt,
		&cat.ID,
		&cat.Name,
		&cat.Description,
	); err != nil {

		c.JSON(http.StatusNotFound, gin.H{
			"error": "Product tidak ditemukan",
		})
		return
	}

	resp := responses.ProductResponse{
		ID:          p.ID,
		CategoryID:  p.CategoryID,
		Name:        p.Name,
		Description: p.Description,
		Price:       p.Price,
		Stock:       p.Stock,
		Image:       helper.FileURL(c, p.Image),
		IsActive:    p.IsActive,
		CreatedAt:   p.CreatedAt,
		UpdatedAt:   p.UpdatedAt,
		Category: responses.ProductCategoryResponse{
			ID:          cat.ID,
			Name:        cat.Name,
			Description: cat.Description,
		},
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Detail Product",
		"data":    resp,
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
