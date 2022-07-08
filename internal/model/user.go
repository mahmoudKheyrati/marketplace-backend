package model

import "time"

type User struct {
	Id             int64     `json:"id"`
	Email          string    `json:"email"`
	Password       string    `json:"password"`
	PhoneNumber    string    `json:"phone_number"`
	FirstName      string    `json:"first_name"`
	LastName       string    `json:"last_name"`
	AvatarUrl      string    `json:"avatar_url"`
	NationalId     string    `json:"national_id"`
	PermissionName string    `json:"permission_name"`
	CreatedAt      time.Time `json:"created_at"`
}
