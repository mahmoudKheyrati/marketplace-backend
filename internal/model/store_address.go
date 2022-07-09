package model

import "time"

type StoreAddress struct {
	Id         int64     `json:"id"`
	StoreId    int64     `json:"store_id"`
	Country    string    `json:"country"`
	Province   string    `json:"province"`
	City       string    `json:"city"`
	Street     string    `json:"street"`
	PostalCode string    `json:"postal_code"`
	CreatedAt  time.Time `json:"created_at"`
}
