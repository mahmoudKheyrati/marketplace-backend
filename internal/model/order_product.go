package model

import "time"

type OrderProduct struct {
	ProductId int64     `json:"product_id"`
	StoreId   int64     `json:"store_id"`
	OrderId   int64     `json:"order_id"`
	Quantity  int64     `json:"quantity"`
	CreatedAt time.Time `json:"created_at"`
}
