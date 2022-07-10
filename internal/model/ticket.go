package model

import "time"

type Ticket struct {
	Id           int64     `json:"id"`
	UserId       int64     `json:"user_id"`
	EmployeeId   *int64    `json:"employee_id"`
	TicketTypeId int64     `json:"ticket_type_id"`
	CreatedAt    time.Time `json:"created_at"`
}
