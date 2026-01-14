package controllers

import (
	"ecommerce-api/config"
	"ecommerce-api/models"
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
	userID := c.GetInt("user_id") // dari JWT

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	offset := (page - 1) * limit

	rows, err := config.DB.Query(`
		SELECT
			t.id, t.qty, t.price, t.total, t.status, t.created_at,
			p.id, p.name
		FROM transactions t
		JOIN products p ON p.id = t.product_id
		WHERE t.buyer_id = $1 OR t.seller_id = $1
		ORDER BY t.created_at DESC
		LIMIT $2 OFFSET $3
	`, userID, limit, offset)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal ambil history"})
		return
	}
	defer rows.Close()

	var data []models.Transaction

	for rows.Next() {
		var trx models.Transaction
		var p models.Product

		rows.Scan(
			&trx.ID,
			&trx.Qty,
			&trx.Price,
			&trx.Total,
			&trx.Status,
			&trx.CreatedAt,
			&p.ID,
			&p.Name,
		)

		trx.Product = &p
		data = append(data, trx)
	}

	var total int
	config.DB.QueryRow(`
		SELECT COUNT(*) FROM transactions
		WHERE buyer_id = $1 OR seller_id = $1
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
		BuyerID   int `json:"buyer_id"`
		SellerID  int `json:"seller_id"`
		Qty       int `json:"qty"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON"})
		return
	}

	var (
		price      float64
		categoryID int
	)

	err := config.DB.QueryRow(`
		SELECT price, category_id
		FROM products
		WHERE id = $1 AND is_active = true
	`, req.ProductID).Scan(&price, &categoryID)

	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Product tidak ditemukan"})
		return
	}

	total := price * float64(req.Qty)

	_, err = config.DB.Exec(`
		INSERT INTO transactions
		(product_id, category_id, buyer_id, seller_id, qty, price, total, status)
		VALUES ($1,$2,$3,$4,$5,$6,$7,'PENDING')
	`,
		req.ProductID,
		categoryID,
		req.BuyerID,
		req.SellerID,
		req.Qty,
		price,
		total,
	)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Gagal membuat transaksi"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Transaksi berhasil dibuat",
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
