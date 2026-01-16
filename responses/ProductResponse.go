package responses

import "time"

type ProductCategoryResponse struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

type ProductResponse struct {
	ID          int                     `json:"id"`
	CategoryID  int                     `json:"category_id"`
	Name        string                  `json:"name"`
	Description string                  `json:"description"`
	Price       float64                 `json:"price"`
	Stock       int                     `json:"stock"`
	Image       string                  `json:"image"`
	IsActive    bool                    `json:"is_active"`
	CreatedAt   time.Time               `json:"created_at"`
	UpdatedAt   *time.Time              `json:"updated_at"`
	Category    ProductCategoryResponse `json:"category"`
}
