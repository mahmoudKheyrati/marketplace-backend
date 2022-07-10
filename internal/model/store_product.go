package model

import "time"

type StoreProduct struct {
	ProductId      int64     `json:"product_id"`
	StoreId        int64     `json:"store_id"`
	OffPercent     float64   `json:"off_percent"`
	MaxOffPrice    float64   `json:"max_off_price"`
	Price          float64   `json:"price"`
	AvailableCount int       `json:"available_count"`
	WarrantyId     *int64    `json:"warranty_id"`
	CreatedAt      time.Time `json:"created_at"`
}
