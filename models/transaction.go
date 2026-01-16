package models

import "time"

type Transaction struct {
	ID        int              `json:"id"`
	BuyerID   int              `json:"buyer_id"`
	SellerID  int              `json:"seller_id"`
	Qty       int              `json:"qty"`
	Price     float64          `json:"price"`
	Total     float64          `json:"total"`
	Status    string           `json:"status"`
	Product   *Product         `json:"product,omitempty"`
	Category  *ProductCategory `json:"product_category,omitempty"`
	Buyer     *User            `json:"buyer,omitempty"`
	Seller    *User            `json:"seller,omitempty"`
	CreatedAt time.Time        `json:"created_at"`
	UpdatedAt *time.Time       `json:"updated_at,omitempty"`
}
