package model

type Category struct {
	Id     int64  `json:"id"`
	Name   string `json:"name"`
	Parent int64  `json:"parent"`
}
