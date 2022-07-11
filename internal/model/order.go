package model

import "time"

type Order struct {
	Id                   int64      `json:"id"`
	Status               string     `json:"status"`
	TrackingCode         *string    `json:"tracking_code"`
	UserId               int64      `json:"user_id"`
	AddressId            *int64     `json:"address_id"`
	ShippingMethodId     *string    `json:"shipping_method_id"`
	AppliedPromotionCode *string    `json:"applied_promotion_code"`
	PayedPrice           *float64   `json:"payed_price"`
	IsPaid               bool       `json:"is_paid"`
	PayDate              *time.Time `json:"pay_date"`
	CreatedAt            time.Time  `json:"created_at"`
}
