package api

import (
	"context"
	"github.com/gofiber/fiber/v2"
	"github.com/mahmoudKheyrati/marketplace-backend/internal/repository"
	"github.com/mahmoudKheyrati/marketplace-backend/pkg"
	"github.com/spf13/cast"
)

type NotificationHandler struct {
	notificationRepo repository.NotificationRepo
}

func NewNotificationHandler(notificationRepo repository.NotificationRepo) *NotificationHandler {
	return &NotificationHandler{notificationRepo: notificationRepo}
}

func (n *NotificationHandler) SubscribeToProductNotification(c *fiber.Ctx) error {
	ctx := context.Background()
	userId := c.Locals(pkg.UserIdKey).(int64)
	productId := cast.ToInt64(c.Params("productId"))

	err := n.notificationRepo.Subscribe(ctx, productId, userId)
	if err != nil {
		pkg.Logger().Error(err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "failed to subscribe on product's available notification"})
	}
	return c.Status(fiber.StatusCreated).JSON(fiber.Map{"status": "ok"})
}

func (n *NotificationHandler) UnSubscribeToProductNotification(c *fiber.Ctx) error {
	ctx := context.Background()
	userId := c.Locals(pkg.UserIdKey).(int64)
	productId := cast.ToInt64(c.Params("productId"))

	err := n.notificationRepo.UnSubscribe(ctx, productId, userId)
	if err != nil {
		pkg.Logger().Error(err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "failed to subscribe on product's available notification"})
	}
	return c.Status(fiber.StatusCreated).JSON(fiber.Map{"status": "ok"})
}

func (n *NotificationHandler) GetAvailableNotifications(c *fiber.Ctx) error {
	ctx := context.Background()
	userId := c.Locals(pkg.UserIdKey).(int64)

	productIds, err := n.notificationRepo.GetAvailableNotifications(ctx, userId)
	if err != nil {
		pkg.Logger().Error(err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "failed to get available notifications"})
	}
	return c.Status(fiber.StatusCreated).JSON(fiber.Map{"status": "ok", "productIds": productIds})

}

func (n *NotificationHandler) SeenNotification(c *fiber.Ctx) error {
	ctx := context.Background()
	userId := c.Locals(pkg.UserIdKey).(int64)
	productId := cast.ToInt64(c.Params("productId"))

	err := n.notificationRepo.MarkNotificationAsSeen(ctx, productId, userId)
	if err != nil {
		pkg.Logger().Error(err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "failed to mark notification as seen"})
	}
	return c.Status(fiber.StatusCreated).JSON(fiber.Map{"status": "ok"})

}
func (n *NotificationHandler) GetPendingNotifications(c *fiber.Ctx) error {
	ctx := context.Background()
	userId := c.Locals(pkg.UserIdKey).(int64)

	productIds, err := n.notificationRepo.GetPendingNotifications(ctx, userId)
	if err != nil {
		pkg.Logger().Error(err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "failed to get pending notifications"})
	}
	return c.Status(fiber.StatusCreated).JSON(fiber.Map{"status": "ok", "productIds": productIds})

}
