package repository

import (
	"context"
	"errors"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/mahmoudKheyrati/marketplace-backend/internal/model"
)

type TicketRepo interface {
	GetAllTickets(ctx context.Context, userId int64) ([]model.Ticket, error)
	GetAllTicketTypes(ctx context.Context) ([]model.TicketType, error)
	CreateTicket(ctx context.Context, ticketTypeId int64, userId int64) (int64, error)
	SendMessageToTicket(ctx context.Context, ticketId int64, senderId int64, text string) error
	GetTicketMessages(ctx context.Context, ticketId int64, userId int64, offset int64) ([]model.TicketMessage, error)

	GetAllUnFinishedTickets(ctx context.Context) ([]model.Ticket, error)
}

type TicketRepoImpl struct {
	db *pgxpool.Pool
}

func NewTicketRepoImpl(db *pgxpool.Pool) *TicketRepoImpl {
	return &TicketRepoImpl{db: db}
}

func (t *TicketRepoImpl) GetAllTickets(ctx context.Context, userId int64) ([]model.Ticket, error) {
	query := `select id, user_id, employee_id, ticket_type_id, created_at
from ticket
where user_id = $1
order by created_at desc`
	rows, err := t.db.Query(ctx, query, userId)
	if err != nil {
		return nil, err
	}
	tickets := make([]model.Ticket, 0)
	for rows.Next() {
		ticket := model.Ticket{}
		err := rows.Scan(
			&ticket.Id,
			&ticket.UserId,
			&ticket.EmployeeId,
			&ticket.TicketTypeId,
			&ticket.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
		tickets = append(tickets, ticket)
	}

	return tickets, nil
}

func (t *TicketRepoImpl) GetAllTicketTypes(ctx context.Context) ([]model.TicketType, error) {
	query := `select id, name, description
from ticket_type
where is_last_version = true`

	rows, err := t.db.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	ticketTypes := make([]model.TicketType, 0)
	for rows.Next() {
		ticketType := model.TicketType{}
		err := rows.Scan(
			&ticketType.Id,
			&ticketType.Name,
			&ticketType.Description,
		)
		if err != nil {
			return nil, err
		}
		ticketTypes = append(ticketTypes, ticketType)
	}

	return ticketTypes, nil
}

func (t *TicketRepoImpl) CreateTicket(ctx context.Context, ticketTypeId int64, userId int64) (int64, error) {
	query := `insert into ticket(user_id, ticket_type_id)
values ($1, $2) returning id`

	row := t.db.QueryRow(ctx, query, userId, ticketTypeId)

	var id int64 = -1
	err := row.Scan(&id)
	return id, err

}

func (t *TicketRepoImpl) SendMessageToTicket(ctx context.Context, ticketId int64, senderId int64, text string) error {
	query := `insert into ticket_message(ticket_id, sender_id, message_text)
values ($1, $2, $3)`
	_, err := t.db.Query(ctx, query, ticketId, senderId, text)
	return err
}

func (t *TicketRepoImpl) GetTicketMessages(ctx context.Context, ticketId int64, userId int64, offset int64) ([]model.TicketMessage, error) {
	query := `select exists( select 1 from ticket where id = $1 and user_id = $2 )`
	rows, err := t.db.Query(ctx, query, ticketId, userId)
	if err != nil {
		return nil, err
	}
	var exists bool
	for rows.Next() {
		err := rows.Scan(&exists)
		if err != nil {
			return nil, err
		}
	}
	if !exists {
		return nil, errors.New("illegal access to tickets.")
	}
	query = `select id, ticket_id, sender_id, message_text, status, created_at
from ticket_message
where ticket_id = $1
order by created_at desc
limit 5 offset $2`
	rows, err = t.db.Query(ctx, query, ticketId, offset)
	if err != nil {
		return nil, err
	}
	ticketMessages := make([]model.TicketMessage, 0)
	for rows.Next() {
		ticketMessage := model.TicketMessage{}
		err := rows.Scan(
			&ticketMessage.Id,
			&ticketMessage.TicketId,
			&ticketMessage.SenderId,
			&ticketMessage.MessageText,
			&ticketMessage.Status,
			&ticketMessage.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
		ticketMessages = append(ticketMessages, ticketMessage)
	}
	return ticketMessages, nil
}

func (t *TicketRepoImpl) GetAllUnFinishedTickets(ctx context.Context) ([]model.Ticket, error) {
	query := `select id, user_id, employee_id, ticket_type_id, is_done, done_at, created_at
from ticket where is_done = false order by created_at desc`
	rows, err := t.db.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	tickets := make([]model.Ticket, 0)
	for rows.Next() {
		ticket := model.Ticket{}
		err := rows.Scan(
			&ticket.Id,
			&ticket.UserId,
			&ticket.EmployeeId,
			&ticket.TicketTypeId,
			&ticket.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
		tickets = append(tickets, ticket)
	}

	return tickets, nil
}
