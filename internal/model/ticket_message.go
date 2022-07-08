package model

import "time"

type TicketMessage struct {
	Id          int64     `json:"id"`
	TicketId    int64     `json:"ticket_id"`
	SenderId    int64     `json:"sender_id"`
	MessageText string    `json:"message_text"`
	Status      string    `json:"status"`
	CreatedAt   time.Time `json:"created_at"`
}
