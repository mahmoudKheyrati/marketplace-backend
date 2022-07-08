package model

import "time"

type Address struct {
	Id              int64
	UserId          int64
	Country         string
	Province        string
	City            string
	Street          string
	PostalCode      string
	HomePhoneNumber string
	CreatedAt       time.Time
}
