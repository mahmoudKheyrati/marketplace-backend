package api

import (
	"context"
	"github.com/gofiber/fiber/v2"
	"github.com/mahmoudKheyrati/marketplace-backend/internal/repository"
	"github.com/mahmoudKheyrati/marketplace-backend/pkg"
	"github.com/spf13/cast"
	"time"
)

type TicketHandler struct {
	ticketRepo repository.TicketRepo
}

func NewTicketHandler(ticketRepo repository.TicketRepo) *TicketHandler {
	return &TicketHandler{ticketRepo: ticketRepo}
}

func (t *TicketHandler) GetAllTickets(c *fiber.Ctx) error {
	ctx := context.Background()
	userId := c.Locals(pkg.UserIdKey).(int64)

	tickets, err := t.ticketRepo.GetAllTickets(ctx, userId)
	if err != nil {
		pkg.Logger().Error(err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": "can not get tickets"})
	}
	return c.JSON(fiber.Map{"tickets": tickets})
}

func (t *TicketHandler) GetAllUnfinishedTickets(c *fiber.Ctx) error {
	ctx := context.Background()
	permissionName := c.Locals(pkg.UserPermissionNameKey).(string)

	if permissionName != "marketplace-admin" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"message": "you don't have access"})
	}
	tickets, err := t.ticketRepo.GetAllUnFinishedTickets(ctx)
	if err != nil {
		pkg.Logger().Error(err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": "can not get tickets"})
	}
	return c.JSON(fiber.Map{"tickets": tickets})
}

func (t *TicketHandler) GetAllTicketTypes(c *fiber.Ctx) error {
	ctx := context.Background()

	ticketTypes, err := t.ticketRepo.GetAllTicketTypes(ctx)
	if err != nil {
		pkg.Logger().Error(err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": "can not get ticket types"})
	}
	return c.JSON(fiber.Map{"ticket_types": ticketTypes})
}

func (t *TicketHandler) CreateTicket(c *fiber.Ctx) error {
	ctx := context.Background()

	userId := c.Locals(pkg.UserIdKey).(int64)
	ticketType := cast.ToInt64(c.Params("ticketTypeId"))

	ticketId, err := t.ticketRepo.CreateTicket(ctx, ticketType, userId)
	if err != nil || ticketId == -1 {
		pkg.Logger().Error(err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": "can not create ticket"})
	}
	return c.JSON(fiber.Map{"ticket_id": ticketId})
}

type SendMessageToTicketRequest struct {
	Text string `json:"text"`
}
type SendMessageToTicketResponse struct {
	TicketId    int64     `json:"ticket_id"`
	SenderId    int64     `json:"sender_id"`
	MessageText string    `json:"message_text"`
	Status      string    `json:"status"`
	CreatedAt   time.Time `json:"created_at"`
}

func (t *TicketHandler) SendMessageToTicket(c *fiber.Ctx) error {
	ctx := context.Background()

	userId := c.Locals(pkg.UserIdKey).(int64)
	ticketId := cast.ToInt64(c.Params("ticketId"))

	var request SendMessageToTicketRequest
	if err := c.BodyParser(&request); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "error", "message": "error in parsing request body"})
	}

	err := t.ticketRepo.SendMessageToTicket(ctx, ticketId, userId, request.Text)
	if err != nil {
		pkg.Logger().Error(err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": "can not send message to ticket"})
	}
	return c.JSON(
		fiber.Map{"message": SendMessageToTicketResponse{
			TicketId:    ticketId,
			SenderId:    userId,
			MessageText: request.Text,
			Status:      "sent",
			CreatedAt:   time.Now(),
		},
			"status": "ok"})
}

func (t *TicketHandler) LoadTicketMessages(c *fiber.Ctx) error {
	ctx := context.Background()

	userId := c.Locals(pkg.UserIdKey).(int64)
	ticketId := cast.ToInt64(c.Params("ticketId"))
	offset := cast.ToInt64(c.Query("offset", "0"))

	messages, err := t.ticketRepo.GetTicketMessages(ctx, ticketId, userId, offset)
	if err != nil {
		pkg.Logger().Error(err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": "can not load ticket messages"})
	}
	return c.JSON(fiber.Map{"ticket_id": ticketId, "next_offset": offset + 5, "messages": messages})
}
