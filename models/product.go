package models

import "time"

type Product struct {
	ID         int `json:"id"`
	CategoryID int `json:"category_id"`

	Name        string `json:"name"`
	Description string `json:"description"`

	Price float64 `json:"price"`
	Stock int     `json:"stock"`

	Image    string `json:"image"`
	IsActive bool   `json:"is_active"`

	Category *ProductCategory `json:"category,omitempty"`

	CreatedAt time.Time  `json:"created_at"`
	CreatedBy *int       `json:"created_by,omitempty"`
	UpdatedAt *time.Time `json:"updated_at,omitempty"`
	UpdatedBy *int       `json:"updated_by,omitempty"`
}
