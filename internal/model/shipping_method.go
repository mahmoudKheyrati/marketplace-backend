package model

import "time"

type ShippingMethod struct {
	Name                       string    `json:"name"`
	ExpectedArrivalWorkingDays int       `json:"expected_arrival_working_days"`
	BaseCost                   int       `json:"base_cost"`
	CreatedAt                  time.Time `json:"created_at"`
}
