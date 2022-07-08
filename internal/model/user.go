package model

import "time"

type User struct {
	Id             int64
	Email          string
	Password       string
	PhoneNumber    string
	FirstName      string
	LastName       string
	AvatarUrl      string
	NationalId     string
	PermissionName string
	CreatedAt      time.Time
}
