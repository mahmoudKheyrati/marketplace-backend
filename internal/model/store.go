package model

import "time"

type Store struct {
	Id          int64     `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	AvatarUrl   string    `json:"avatar_url"`
	Owner       int64     `json:"owner"`
	Creator     int64     `json:"creator"`
	CreatedAt   time.Time `json:"created_at"`
}
