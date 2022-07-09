package model

import "time"

type Review struct {
	Id         int64     `json:"id"`
	ProductId  int64     `json:"product_id"`
	StoreId    int64     `json:"store_id"`
	UserId     int64     `json:"user_id"`
	Rate       float64   `json:"rate"`
	ReviewText string    `json:"review_text"`
	CreatedAt  time.Time `json:"created_at"`
	UpVotes    int64     `json:"up_votes"`
	DownVotes  int64     `json:"down_votes"`
}
