package controllers

import (
	"ecommerce-api/config"
	helper "ecommerce-api/helpers"
	"ecommerce-api/models"
	"ecommerce-api/responses"
	"fmt"
	"math"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func GetTransactions(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))

	if page < 1 {
		page = 1
	}
	offset := (page - 1) * limit

	rows, err := config.DB.Query(`
		SELECT
			t.id, t.qty, t.price, t.total, t.status, t.created_at,
			p.id, p.name,
			c.id, c.name
		FROM transactions t
		JOIN products p ON p.id = t.product_id
		JOIN categories c ON c.id = t.category_id
		ORDER BY t.created_at DESC
		LIMIT $1 OFFSET $2
	`, limit, offset)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal ambil data"})
		return
	}
	defer rows.Close()

	var data []models.Transaction

	for rows.Next() {
		var trx models.Transaction
		var p models.Product
		var cat models.ProductCategory

		rows.Scan(
			&trx.ID,
			&trx.Qty,
			&trx.Price,
			&trx.Total,
			&trx.Status,
			&trx.CreatedAt,
			&p.ID,
			&p.Name,
			&cat.ID,
			&cat.Name,
		)

		trx.Product = &p
		trx.Category = &cat
		data = append(data, trx)
	}

	var total int
	config.DB.QueryRow(`SELECT COUNT(*) FROM transactions`).Scan(&total)

	c.JSON(http.StatusOK, gin.H{
		"message": "List Transactions",
		"data":    data,
		"meta": gin.H{
			"page":       page,
			"limit":      limit,
			"total":      total,
			"total_page": int(math.Ceil(float64(total) / float64(limit))),
		},
	})
}

func TransactionHistory(c *gin.Context) {
	userID := c.GetInt("user_id")

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	offset := (page - 1) * limit

	rows, err := config.DB.Query(`
		SELECT
			t.id, t.qty, t.price, t.total, t.status, t.created_at,
			p.id, p.name, p.image
		FROM transactions t
		JOIN products p ON p.id = t.product_id
		WHERE t.buyer_id = $1
		ORDER BY t.created_at DESC
		LIMIT $2 OFFSET $3
	`, userID, limit, offset)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal ambil history"})
		return
	}
	defer rows.Close()

	var data []responses.TransactionHistoryResponse

	for rows.Next() {
		var trx responses.TransactionHistoryResponse

		err := rows.Scan(
			&trx.ID,
			&trx.Qty,
			&trx.Price,
			&trx.Total,
			&trx.Status,
			&trx.CreatedAt,
			&trx.Product.ID,
			&trx.Product.Name,
			&trx.Product.Image,
		)

		trx.Product.Image = helper.FileURL(c, trx.Product.Image)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Scan error"})
			return
		}

		data = append(data, trx)
	}

	var total int
	config.DB.QueryRow(`
		SELECT COUNT(*) FROM transactions WHERE buyer_id = $1
	`, userID).Scan(&total)

	c.JSON(http.StatusOK, gin.H{
		"message": "Transaction History",
		"data":    data,
		"meta": gin.H{
			"page":       page,
			"limit":      limit,
			"total":      total,
			"total_page": int(math.Ceil(float64(total) / float64(limit))),
		},
	})
}

func CreateTransaction(c *gin.Context) {
	var req struct {
		ProductID int `json:"product_id"`
		SellerID  int `json:"seller_id"`
		Qty       int `json:"qty"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON"})
		return
	}

	if req.Qty <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Qty minimal 1"})
		return
	}

	buyerID := c.GetInt("user_id")

	var (
		price      float64
		stock      int
		categoryID int
	)

	err := config.DB.QueryRow(`
		SELECT price, stock, category_id
		FROM products
		WHERE id = $1 AND is_active = true
	`, req.ProductID).Scan(&price, &stock, &categoryID)

	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Product tidak ditemukan"})
		return
	}

	if stock < req.Qty {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Stock tidak mencukupi"})
		return
	}

	total := price * float64(req.Qty)

	tx, err := config.DB.Begin()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mulai transaksi"})
		return
	}

	defer tx.Rollback()

	_, err = tx.Exec(`
		INSERT INTO transactions (buyer_id, seller_id, product_id, qty, price, total, category_id, status)
		VALUES ($1,$2,$3,$4,$5,$6,$7,'PENDING')
	`, buyerID, req.SellerID, req.ProductID, req.Qty, price, total, categoryID)

	if err != nil {
		fmt.Println("SQL ERROR:", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Gagal membuat transaksi"})
		return
	}

	_, err = tx.Exec(`
		UPDATE products
		SET stock = stock - $1
		WHERE id = $2
	`, req.Qty, req.ProductID)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal update stock"})
		return
	}

	tx.Commit()

	c.JSON(http.StatusCreated, gin.H{
		"message": "Transaksi berhasil",
	})
}

func GetTransactionByID(c *gin.Context) {
	id := c.Param("id")

	var trx models.Transaction
	var product models.Product
	var category models.ProductCategory
	var buyer models.User
	var seller models.User

	err := config.DB.QueryRow(`
		SELECT 
			t.id, t.qty, t.price, t.total, t.status, t.created_at,

			p.id, p.name,
			c.id, c.name,

			b.id, b.name,
			s.id, s.name
		FROM transactions t
		JOIN products p ON p.id = t.product_id
		JOIN categories c ON c.id = t.category_id
		JOIN users b ON b.id = t.buyer_id
		JOIN users s ON s.id = t.seller_id
		WHERE t.id = $1
	`, id).Scan(
		&trx.ID,
		&trx.Qty,
		&trx.Price,
		&trx.Total,
		&trx.Status,
		&trx.CreatedAt,
		&product.ID,
		&product.Name,
		&category.ID,
		&category.Name,
		&buyer.ID,
		&buyer.Name,
		&seller.ID,
		&seller.Name,
	)

	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Transaksi tidak ditemukan"})
		return
	}

	trx.Product = &product
	trx.Category = &category
	trx.Buyer = &buyer
	trx.Seller = &seller

	c.JSON(http.StatusOK, gin.H{
		"message": "Detail Transaction",
		"data":    trx,
	})
}
