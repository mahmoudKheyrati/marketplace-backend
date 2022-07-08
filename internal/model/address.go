package model

import "time"

type Address struct {
	Id              int64     `json:"id"`
	UserId          int64     `json:"user_id"`
	Country         string    `json:"country"`
	Province        string    `json:"province"`
	City            string    `json:"city"`
	Street          string    `json:"street"`
	PostalCode      string    `json:"postal_code"`
	HomePhoneNumber string    `json:"home_phone_number"`
	CreatedAt       time.Time `json:"created_at"`
}
