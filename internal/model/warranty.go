package model

import "time"

type Warranty struct {
	Id           int64     `json:"id"`
	Name         string    `json:"name"`
	WarrantyType string    `json:"warranty_type"`
	Month        int       `json:"month"`
	CreatedAt    time.Time `json:"created_at"`
}
