package responses

import "time"

type TransactionHistoryResponse struct {
	ID        int                 `json:"id"`
	Qty       int                 `json:"qty"`
	Price     float64             `json:"price"`
	Total     float64             `json:"total"`
	Status    string              `json:"status"`
	CreatedAt time.Time           `json:"created_at"`
	Product   ProductMiniResponse `json:"product"`
}

type ProductMiniResponse struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Image string `json:"image"`
}
